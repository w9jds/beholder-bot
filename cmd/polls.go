package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func updatePollAnswers(session *discordgo.Session, reactionAdd *discordgo.MessageReactionAdd) {
	message, err := session.ChannelMessage(reactionAdd.ChannelID, reactionAdd.MessageID)
	if err != nil {
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		log.Println(err)
		return
	}

	if len(message.Embeds) == 0 || message.Embeds[0].Title != dayPollTitle {
		return
	}

	user, err := session.GuildMember(reactionAdd.GuildID, reactionAdd.UserID)
	if err != nil {
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		log.Println(err)
		return
	}

	// No Emoji ID is returned right now, so dirty hack is to use the name for now.
	switch reactionAdd.Emoji.Name {
	case "🇲":
		message.Embeds[0].Fields[0].Value = addUser(message.Embeds[0].Fields[0].Value, user)
		break
	case "🇹":
		message.Embeds[0].Fields[1].Value = addUser(message.Embeds[0].Fields[1].Value, user)
		break
	case "🇼":
		message.Embeds[0].Fields[2].Value = addUser(message.Embeds[0].Fields[2].Value, user)
		break
	case "🇭":
		message.Embeds[0].Fields[3].Value = addUser(message.Embeds[0].Fields[3].Value, user)
		break
	case "🇫":
		message.Embeds[0].Fields[4].Value = addUser(message.Embeds[0].Fields[4].Value, user)
		break
	case "🇸":
		message.Embeds[0].Fields[5].Value = addUser(message.Embeds[0].Fields[5].Value, user)
		break
	case "🇺":
		message.Embeds[0].Fields[6].Value = addUser(message.Embeds[0].Fields[6].Value, user)
		break
	default:
		removeReaction(session, reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.ID, reactionAdd.UserID)
		return
	}

	_, err = session.ChannelMessageEditEmbed(reactionAdd.ChannelID, message.ID, message.Embeds[0])
	if err != nil {
		fmt.Println(err)
	}
}

func addUser(value string, user *discordgo.Member) string {
	var users []string
	if value == "None" {
		users = []string{}
	} else {
		users = strings.Split(value, ", ")
	}

	if strings.Index(value, user.Mention()) > -1 {
		return value
	}

	return strings.Join(append(users, user.Mention()), ", ")
}

func removeReaction(session *discordgo.Session, channelID, messageID, emojiID, userID string) {
	err := session.MessageReactionRemove(channelID, messageID, emojiID, userID)
	if err != nil {
		log.Println(err)
	}
}

func pollBestDay(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !isGameDM(message.GuildID, message.ChannelID, message.Author.ID) {
		return
	}

	_, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "@here The DM has initiated a poll to see what day best works for all players, add a reaction for the days of the week that you are available to play!",
		Embed: &discordgo.MessageEmbed{
			Title: dayPollTitle,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   ":regional_indicator_m:onday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   ":regional_indicator_t:uesday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   ":regional_indicator_w:ednesday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "T:regional_indicator_h:ursday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   ":regional_indicator_f:riday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   ":regional_indicator_s:aturday",
					Value:  "None",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "S:regional_indicator_u:nday",
					Value:  "None",
					Inline: false,
				},
			},
		},
	})

	if err != nil {
		log.Println(err)
	}
}
