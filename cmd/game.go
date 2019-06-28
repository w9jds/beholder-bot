package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	dayPollTitle string = "Best day of the week poll"
)

func createNewGame(session *discordgo.Session, message *discordgo.MessageCreate) {
	overwrites := []*discordgo.PermissionOverwrite{
		getDmPermissions(message.Author),
		&discordgo.PermissionOverwrite{
			ID:    message.GuildID,
			Type:  "role",
			Allow: 0,
			Deny:  discordgo.PermissionAll,
		},
	}

	var name string
	if len(message.Mentions) > 0 {
		name = strings.Trim(message.Content[strings.Index(message.Content, " "):strings.Index(message.Content, "@")-1], " ")
	} else {
		words := strings.Split(message.Content, " ")
		name = strings.Join(words[1:], " ")
	}

	for _, user := range message.Mentions {
		overwrites = append(overwrites, getPlayerPermissions(user))
	}

	channel, err := session.GuildChannelCreateComplex(message.GuildID, discordgo.GuildChannelCreateData{
		Name:                 name,
		Type:                 discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: overwrites,
	})
	if err != nil {
		log.Println("Error creating category for your party")
		return
	}

	textChannel, err := session.GuildChannelCreateComplex(message.GuildID, discordgo.GuildChannelCreateData{
		ParentID: channel.ID,
		Name:     strings.ReplaceAll(name, " ", "-"),
		Type:     discordgo.ChannelTypeGuildText,
	})
	if err != nil {
		log.Println("Error creating text chat for your party")
		return
	}

	voiceChannel, err := session.GuildChannelCreateComplex(message.GuildID, discordgo.GuildChannelCreateData{
		ParentID: channel.ID,
		Name:     name,
		Type:     discordgo.ChannelTypeGuildVoice,
	})
	if err != nil {
		log.Println("Error creating voice chat for your party")
		return
	}

	storeNewGame(message.GuildID, channel.ID, textChannel.ID, voiceChannel.ID, message.Author.ID)
}

func setNextSession(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !isGameDM(message.GuildID, message.ChannelID, message.Author.ID) {
		return
	}

}

func getDmPermissions(user *discordgo.User) *discordgo.PermissionOverwrite {
	return &discordgo.PermissionOverwrite{
		ID:   user.ID,
		Type: "member",
		Allow: discordgo.PermissionAllText |
			discordgo.PermissionAllVoice |
			discordgo.PermissionAddReactions |
			0x00000100,
		Deny: 0,
	}
}

func getPlayerPermissions(user *discordgo.User) *discordgo.PermissionOverwrite {
	return &discordgo.PermissionOverwrite{
		ID:   user.ID,
		Type: "member",
		Allow: discordgo.PermissionEmbedLinks |
			discordgo.PermissionAddReactions |
			discordgo.PermissionReadMessages |
			discordgo.PermissionReadMessageHistory |
			discordgo.PermissionSendMessages |
			discordgo.PermissionSendTTSMessages |
			discordgo.PermissionAttachFiles |
			discordgo.PermissionVoiceConnect |
			discordgo.PermissionVoiceSpeak |
			discordgo.PermissionVoiceUseVAD,
		Deny: discordgo.PermissionVoiceMuteMembers |
			discordgo.PermissionVoiceDeafenMembers |
			discordgo.PermissionVoiceMoveMembers |
			discordgo.PermissionManageMessages,
	}
}
