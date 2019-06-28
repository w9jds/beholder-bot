package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/bwmarrin/discordgo"
)

var (
	httpClient *http.Client
	discord    *discordgo.Session
	postgres   *sql.DB

	botToken string
	dsn      string
)

func setupEnv() bool {
	botToken = strings.Trim(os.Getenv("BOT_TOKEN"), " ")
	if botToken == "" {
		log.Fatal("Environment variable 'BOT_TOKEN' requires a value to be able to push messages to discord")
	}

	host := strings.Trim(os.Getenv("POSTGRES_HOST"), " ")
	if host == "" {
		log.Fatal("Environment variable `POSTGRES_HOST` is required to connect to the systems database")
	}

	user := strings.Trim(os.Getenv("POSTGRES_USER"), " ")
	if user == "" {
		log.Fatal("Environment variable `POSTGRES_USER` is required to connect to the systems database")
	}

	dbname := strings.Trim(os.Getenv("POSTGRES_DB"), " ")
	if dbname == "" {
		log.Fatal("Environment variable `POSTGRES_DB` is required to connect to the systems database")
	}

	password := strings.Trim(os.Getenv("POSTGRES_PASSWORD"), " ")
	if password == "" {
		log.Fatal("Environment variable `POSTGRES_PASSWORD` is required to connect to the systems database")
	}

	dsn = fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, password)

	return true
}

func main() {
	var err error
	isReady := setupEnv()

	if isReady == true {
		postgres, err = sql.Open("cloudsqlpostgres", dsn)
		if err != nil {
			log.Panic(err)
		}

		defer postgres.Close()

		contents, err := ioutil.ReadFile("../setup.sql")
		if err != nil {
			log.Panic(err)
		}

		_, err = postgres.Exec(string(contents))
		if err != nil {
			log.Panic(err)
		}

		discord, err = discordgo.New("Bot " + botToken)
		if err != nil {
			log.Panic("discordgo: ", err)
		}

		discord.AddHandler(ready)
		discord.AddHandler(messageCreate)
		discord.AddHandler(messageReactionAdd)
		discord.AddHandler(messageDeleted)
		discord.AddHandler(channelDeleted)

		err = discord.Open()
		if err != nil {
			log.Panic("discordgo: ", err)
		}

		defer discord.Close()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func ready(session *discordgo.Session, ready *discordgo.Ready) {
	log.Println("Beholder has started! All systems green.")
}

func messageReactionAdd(session *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	if isGameChannel(reactionAdd.GuildID, reactionAdd.ChannelID) {
		go updatePollAnswers(session, reactionAdd)
	}
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!createadventureparty") {
		go createNewGame(session, message)
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!setnextsession") {
		go setNextSession(session, message)
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!pollbestday") {
		go pollBestDay(session, message)
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!addmap") {
		go addNewMap(session, message)
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!getmap") {
		go getMap(session, message)
	}
}

func messageDeleted(session *discordgo.Session, message *discordgo.MessageDelete) {
	go removeMaps(message.GuildID, message.ChannelID, message.ID)
}

func channelDeleted(session *discordgo.Session, channel *discordgo.ChannelDelete) {
	go removeChannel(channel.GuildID, channel.ID)
}
