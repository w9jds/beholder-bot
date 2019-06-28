package main

import (
	"database/sql"
	"fmt"
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
	var error error
	isReady := setupEnv()

	if isReady == true {
		postgres, error = sql.Open("cloudsqlpostgres", dsn)
		if error != nil {
			log.Fatal(error)
		}

		defer postgres.Close()

		discord, error = discordgo.New("Bot " + botToken)
		if error != nil {
			log.Fatal("discordgo: ", error)
		}

		discord.AddHandler(ready)
		discord.AddHandler(messageCreate)
		discord.AddHandler(messageReactionAdd)

		error = discord.Open()
		if error != nil {
			log.Fatal("discordgo: ", error)
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
	if isGameChannel(reactionAdd.GuildID, reactionAdd.ChannelID) && !isGameDM(reactionAdd.GuildID, reactionAdd.ChannelID, reactionAdd.UserID) {
		go updatePollAnswers(session, reactionAdd)
	}
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if strings.HasPrefix(strings.ToLower(message.Content), "!createdndparty") {
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

}
