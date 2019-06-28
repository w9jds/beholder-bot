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
	name := strings.Trim(message.Content[strings.Index(message.Content, " "):strings.Index(message.Content, "@")-1], " ")
	overwrites := []*discordgo.PermissionOverwrite{
		getDmPermissions(message.Author),
		&discordgo.PermissionOverwrite{
			ID:    message.GuildID,
			Type:  "role",
			Allow: 0,
			Deny:  discordgo.PermissionAll,
		},
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

func addNewMap(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !isGameDM(message.GuildID, message.ChannelID, message.Author.ID) {
		return
	}
}

func updatePollAnswers(session *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	message, err := session.ChannelMessage(reactionAdd.ChannelID, reactionAdd.MessageID)
	if err != nil {
		log.Fatal(err)
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		return
	}

	if len(message.Embeds) == 0 || message.Embeds[0].Title != dayPollTitle {
		return
	}

	user, err := session.GuildMember(reactionAdd.GuildID, reactionAdd.UserID)
	if err != nil {
		log.Fatal(err)
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		return
	}

	switch reactionAdd.Emoji.Name {
	case "regional_indicator_m":
		message.Embeds[0].Fields[0].Value = addUser(message.Embeds[0].Fields[0].Value, user)
		break
	case "regional_indicator_t":
		message.Embeds[0].Fields[1].Value = addUser(message.Embeds[0].Fields[1].Value, user)
		break
	case "regional_indicator_w":
		message.Embeds[0].Fields[2].Value = addUser(message.Embeds[0].Fields[2].Value, user)
		break
	case "regional_indicator_h":
		message.Embeds[0].Fields[3].Value = addUser(message.Embeds[0].Fields[3].Value, user)
		break
	case "regional_indicator_f":
		message.Embeds[0].Fields[4].Value = addUser(message.Embeds[0].Fields[4].Value, user)
		break
	case "regional_indicator_s":
		message.Embeds[0].Fields[5].Value = addUser(message.Embeds[0].Fields[5].Value, user)
		break
	case "regional_indicator_u":
		message.Embeds[0].Fields[6].Value = addUser(message.Embeds[0].Fields[6].Value, user)
		break
	default:
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		return
	}

	session.ChannelMessageEditEmbed(reactionAdd.ChannelID, message.ChannelID, message.Embeds[0])
}

func addUser(value string, user *discordgo.Member) string {
	users := strings.Split(value, ", ")

	if strings.Index(value, user.Mention()) > -1 {
		return value
	}

	return strings.Join(append(users, user.Mention()), ", ")
}

func removeReaction(session *discordgo.Session, channelID string, messageID string, emojiID string, userID string) {
	err := session.MessageReactionRemove(channelID, messageID, emojiID, userID)
	if err != nil {
		log.Fatal(err)
	}
}

func pollBestDay(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !isGameDM(message.GuildID, message.ChannelID, message.Author.ID) {
		return
	}

	session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Title:       dayPollTitle,
		Description: "The DM has initiated a poll to see what day best works for all players, add a reaction for the days of the week that you are available.",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  ":regional_indicator_m:onday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  ":regional_indicator_t:uesday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  ":regional_indicator_w:ednesday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  "T:regional_indicator_h:ursday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  ":regional_indicator_f:riday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  ":regional_indicator_s:aturday",
				Value: "",
			},
			&discordgo.MessageEmbedField{
				Name:  "S:regional_indicator_u:nday",
				Value: "",
			},
		},
	})
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
