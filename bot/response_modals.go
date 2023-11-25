package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

func (bot *ModeratorBot) RespondToModeratePositivelyModalFromUserContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: globals.ModerationSuccessful,
					Color: globals.Purple,
				},
			},
		},
	})
}

func (bot *ModeratorBot) RespondToModerateNegativelyModalFromUserContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: globals.ModerationSuccessful,
					Color: globals.Purple,
				},
			},
		},
	})
}

func (bot *ModeratorBot) RespondToModeratePositivelyModalFromMessageContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: globals.ModerationSuccessful,
					Color: globals.Purple,
				},
			},
		},
	})
}

func (bot *ModeratorBot) RespondToModerateNegativelyModalFromMessageContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: globals.ModerationSuccessful,
					Color: globals.Purple,
				},
			},
		},
	})
}

// TODO: align with SaveEvidenceFromMessageContext when in a good state
// Modal command, copies message details to a configured channel
func (bot *ModeratorBot) SaveEvidenceFromModalSubmission(i *discordgo.InteractionCreate) {
	data := i.Interaction.ModalSubmitData()

	// message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	// if message.ID == "" {
	// 	log.Warn("no message was provided")
	// }

	// fields := []*discordgo.MessageEmbedField{
	// 	{
	// 		Name:   "Original message timestamp",
	// 		Value:  fmt.Sprintf("%s (<t:%v:R>)", message.Timestamp.Format(time.RFC1123Z), message.Timestamp.Unix()),
	// 		Inline: false,
	// 	},
	// 	{
	// 		Name:   "Author of message",
	// 		Value:  fmt.Sprintf("<@%s>", message.Author.ID),
	// 		Inline: true,
	// 	},
	// 	{
	// 		Name:   "Link to original message",
	// 		Value:  fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, message.ChannelID, message.ID),
	// 		Inline: true,
	// 	},
	// 	{
	// 		Name:  "Content of original message",
	// 		Value: message.Content,
	// 	},
	// }

	// if len(message.Attachments) > 0 {
	// 	fields = append(fields, &discordgo.MessageEmbedField{
	// 		Name:  "Attachments",
	// 		Value: fmt.Sprintf("%v", len(message.Attachments)),
	// 	})
	// 	for _, attachment := range message.Attachments {
	// 		fields = append(fields, &discordgo.MessageEmbedField{
	// 			Name:  attachment.Filename,
	// 			Value: attachment.URL,
	// 		})
	// 	}
	// }

	// // TODO: more information
	// ms := discordgo.MessageSend{
	// 	Embeds: []*discordgo.MessageEmbed{{
	// 		Title: "Moder8s Evidence",
	// 		Description: fmt.Sprintf("Collected by <@%s> at %s \\n(<t:%v:R>)",
	// 			i.Interaction.Member.User.ID,
	// 			time.Now().Format(time.RFC1123Z),
	// 			time.Now().Unix(),
	// 		),
	// 		Fields: fields,
	// 		Color:  globals.Purple,
	// 	}},
	// }

	// // TODO: save event info
	// sc := bot.getServerConfig(i.GuildID)
	// _, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	// if err != nil {
	// 	log.Warn("Unable to send message %w", err)
	// }

	// TODO: change this to messagesend
	// _ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: "Evidence saved to <#" + sc.EvidenceChannelSettingID + ">",
	// 		Flags:   discordgo.MessageFlagsEphemeral,
	// 	},
	// })

	userID := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	channelID := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu).CustomID
	messageID := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu).CustomID
	reason := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu).CustomID

	sc := bot.getServerConfig(i.GuildID)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "reason was " + reason,
		},
	})

	ms := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title: "Title goes here",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "User",
					Value: fmt.Sprintf("<@%s>", userID),
				},
				{
					Name:  "Link to original message",
					Value: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", sc.DiscordId, channelID, messageID),
				},
			},
		}},
	}

	_, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("Unable to send message %w", err)
	}
}

func (bot *ModeratorBot) ShowLowReputationModal(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.NotifyWhenReputationIsBelowSetting,
			Title:    "What should the new value be?",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "12345",
						Label:       "Reputation",
						Style:       discordgo.TextInputShort,
						Placeholder: "this user is amazing",
						Required:    true,
						MinLength:   1,
						MaxLength:   4,
					},
				},
			}},
		},
	})
}

func (bot *ModeratorBot) ShowHighReputationModal(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.NotifyWhenReputationIsAboveSetting,
			Title:    "What should the new value be?",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "12345",
						Label:       "Reputation",
						Style:       discordgo.TextInputShort,
						Placeholder: "this user is amazing",
						Required:    true,
						MinLength:   1,
						MaxLength:   4,
					},
				},
			}},
		},
	})
}
