package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tyzbit/go-discord-modtools/globals"
)

func (bot *ModeratorBot) DocumentBehaviorFromMessage(i *discordgo.InteractionCreate, message *discordgo.Message) (resp *discordgo.InteractionResponse) {
	user := bot.GetModeratedUser(i.GuildID, message.Author.ID)
	var fields []*discordgo.MessageEmbedField
	var messageType discordgo.InteractionResponseType
	var authorID string
	if len(message.Embeds) > 0 {
		fields = message.Embeds[0].Fields
		messageType = discordgo.InteractionResponseUpdateMessage
		authorID = getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[1].Value)
		for idx, field := range fields {
			if field.Name == globals.CurrentReputation {
				user := bot.GetModeratedUser(i.GuildID, authorID)
				fields[idx].Value = fmt.Sprintf("%v", user.Reputation.Int64)
			}
		}
	} else {
		messageType = discordgo.InteractionResponseChannelMessageWithSource
		authorID = message.Author.ID
		fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Original message timestamp",
				Value:  fmt.Sprintf("%s (<t:%v:R>)", message.Timestamp.Format(time.RFC1123Z), message.Timestamp.Unix()),
				Inline: false,
			},
			{
				Name:   "Author of message",
				Value:  fmt.Sprintf("<@%s>", user.UserID),
				Inline: true,
			},
			{
				Name:   "Initial reputation",
				Value:  fmt.Sprintf("%v", user.Reputation.Int64),
				Inline: true,
			},
			{
				Name:   globals.CurrentReputation,
				Value:  fmt.Sprintf("%v", user.Reputation.Int64),
				Inline: true,
			},
			{
				Name:   "Link to original message",
				Value:  fmt.Sprintf(globals.MessageURLTemplate, i.Interaction.GuildID, message.ChannelID, message.ID),
				Inline: true,
			},
			{
				Name:  globals.OriginalMessageContent,
				Value: message.Content,
			},
		}

		if len(message.Attachments) > 0 {
			attachmentList := ""
			for _, attachment := range message.Attachments {
				attachmentList = attachmentList + fmt.Sprintf("[%s](%v)\n", attachment.Filename, attachment.URL)
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf(globals.Attachments+" (%v)", len(message.Attachments)),
				Value: attachmentList,
			})
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "Collected by",
			Value: fmt.Sprintf(`<@%s>
									%s (<t:%v:R>)`,
				i.Interaction.Member.User.ID,
				time.Now().Format(time.RFC1123Z),
				time.Now().Unix()),
		})
	}

	return &discordgo.InteractionResponse{
		Type: messageType,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Evidence report",
					Description: fmt.Sprintf("Document user behavior for for <@%v> - good, bad, or noteworthy", authorID),
					Color:       globals.Purple,
					Fields:      fields,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: globals.DecreaseUserReputation,
						Label:    globals.DecreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.IncreaseUserReputation,
						Label:    globals.IncreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.TakeEvidenceNotes,
						Label:    globals.TakeEvidenceNotes,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.SubmitReport,
						Label:    globals.SubmitReport,
						Style:    discordgo.PrimaryButton,
					},
				},
			}},
		},
	}
}
