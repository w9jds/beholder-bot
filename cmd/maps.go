package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func addNewMap(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !isGameDM(message.GuildID, message.ChannelID, message.Author.ID) {
		return
	}

	if len(message.Attachments) < 1 {
		session.ChannelMessageSend(message.ChannelID, "A map must contain a name and attachment!")
		return
	}

	words := strings.Split(message.Content, " ")
	if !storeNewMap(message.GuildID, message.ChannelID, strings.Join(words[1:], " "), message.ID) {
		session.ChannelMessageSend(message.ChannelID, "Unable to store the map, please try again.")
	}
}

func getMap(session *discordgo.Session, message *discordgo.MessageCreate) {
	words := strings.Split(message.Content, " ")
	name := strings.Join(words[1:], " ")

	messageID, err := getStoredMap(message.GuildID, message.ChannelID, name)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("A map by the name of `%s` was not found, please try again.", name))
		return
	}

	originalMessage, err := session.ChannelMessage(message.ChannelID, messageID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "The message containing the map couldn't be found.")
		return
	}

	session.ChannelMessageSend(message.ChannelID, originalMessage.Attachments[0].ProxyURL)
}
