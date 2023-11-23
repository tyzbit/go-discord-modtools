package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// TODO: I think all of these need to log events

// Moderation modal menu
func (bot *ModeratorBot) ModerateMenuFromMessage(i *discordgo.InteractionCreate) {
	message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if message.ID == "" {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ModerateUsingMessage,
			Title:    "Moderate user",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "12345",
						Label:       "Reputation",
						Style:       discordgo.TextInputShort,
						Placeholder: "this user is amazing",
						Required:    true,
						MinLength:   1,
						MaxLength:   300,
					},
				},
			}},
		},
	})

	// _ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		CustomID: globals.ModerateModal,
	// 		Title:    "Moderate " + i.Member.User.Username,
	// 		Embeds: []*discordgo.MessageEmbed{{
	// 			Title: "Embed!",
	// 			Fields: []*discordgo.MessageEmbedField{{
	// 				Name:  "Field1!",
	// 				Value: "Value!",
	// 			}},
	// 		}},
	// 		Components: []discordgo.MessageComponent{
	// 			discordgo.ActionsRow{
	// 				Components: []discordgo.MessageComponent{
	// 					discordgo.SelectMenu{
	// 						Placeholder: "User",
	// 						MenuType:    discordgo.UserSelectMenu,
	// 						CustomID:    i.Member.User.ID,
	// 						Options: []discordgo.SelectMenuOption{{
	// 							Label: "UserID",
	// 							Value: "Saved",
	// 						}},
	// 					},
	// 				},
	// 			},
	// 			discordgo.ActionsRow{
	// 				Components: []discordgo.MessageComponent{
	// 					discordgo.SelectMenu{
	// 						Placeholder:  "Channel",
	// 						MenuType:     discordgo.ChannelSelectMenu,
	// 						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
	// 						CustomID:     message.ChannelID,
	// 						Options: []discordgo.SelectMenuOption{{
	// 							Label: "Channel",
	// 							Value: "Saved",
	// 						}},
	// 					},
	// 				},
	// 			},
	// 			discordgo.ActionsRow{
	// 				Components: []discordgo.MessageComponent{
	// 					discordgo.SelectMenu{
	// 						Placeholder: "Message",
	// 						CustomID:    message.ID,
	// 						Options: []discordgo.SelectMenuOption{{
	// 							Label: "Message",
	// 							Value: "Saved",
	// 						}},
	// 					},
	// 				},
	// 			},
	// 			discordgo.ActionsRow{
	// 				Components: []discordgo.MessageComponent{
	// 					discordgo.TextInput{
	// 						CustomID:    globals.ReasonOption,
	// 						Label:       "Reason",
	// 						Style:       discordgo.TextInputShort,
	// 						Placeholder: "Why are you moderating this user or content?",
	// 						Required:    true,
	// 						MinLength:   1,
	// 						MaxLength:   500,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// })

}

func (bot *ModeratorBot) ModerateMenuFromUser(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User == nil {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ModerateModal,
			Title:    "Moderate " + i.Member.User.Username,
			Embeds: []*discordgo.MessageEmbed{{
				Title: "Embed!",
				Fields: []*discordgo.MessageEmbedField{{
					Name:  "Field1!",
					Value: "Value!",
				}},
			}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder: "User",
							CustomID:    i.Member.User.ID,
							Options: []discordgo.SelectMenuOption{{
								Label: "UserID",
								Value: "Saved",
							}},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.ReasonOption,
							Label:       "Reason",
							Style:       discordgo.TextInputShort,
							Placeholder: "Why are you moderating this user or content?",
							Required:    true,
							MinLength:   1,
							MaxLength:   500,
						},
					},
				},
			},
		},
	})
}

func (bot *ModeratorBot) ModerateActionFromMessage(i *discordgo.InteractionCreate) {
	bot.CollectMessageAsEvidenceFromMessageModal(i)
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Moderation action here",
			Embeds: []*discordgo.MessageEmbed{{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Moderator Username",
						Value:  i.Interaction.Member.User.Username,
						Inline: true,
					},
					{
						Name:   "Channel",
						Value:  fmt.Sprintf("<#" + i.Interaction.ChannelID + ">"),
						Inline: true,
					},
				},
			}},
		},
	})
}

func (bot *ModeratorBot) ModerateActionFromUser(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Moderation action here",
			Embeds: []*discordgo.MessageEmbed{{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Moderator Username",
						Value:  i.Interaction.Member.User.Username,
						Inline: true,
					},
					{
						Name:   "Channel",
						Value:  fmt.Sprintf("<#" + i.Interaction.ChannelID + ">"),
						Inline: true,
					},
				},
			}},
		},
	})
}

// App command, copies message details to a configured channel
func (bot *ModeratorBot) CollectMessageAsEvidence(i *discordgo.InteractionCreate) {
	message := i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	sc := bot.getServerConfig(i.GuildID)
	if message == nil {
		log.Warn("message was not provided")
	}

	ms := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title: "Title goes here",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "User",
					Value: fmt.Sprintf("<@%s>", message.Author.ID),
				},
				{
					Name:  "Link to original message",
					Value: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", message.GuildID, message.ChannelID, message.ID),
				},
			},
		}},
	}

	message, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("Unable to send message %w", err)
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message.Content,
			Embeds: []*discordgo.MessageEmbed{{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Username",
						Value:  message.Author.Username,
						Inline: true,
					},
					{
						Name:   "Channel",
						Value:  fmt.Sprintf("<#" + message.ChannelID + ">"),
						Inline: true,
					},
					{
						Name:   "Timestamp",
						Value:  message.Timestamp.Format(time.RFC1123Z),
						Inline: false,
					},
				},
			}},
			TTS:        message.TTS,
			Components: message.Components,
			//Files: m.Attachments,
			// AllowedMentions,
			//File: ,
			// Embed: m.Embeds[],

		},
	})
}

// Modal command, copies message details to a configured channel
func (bot *ModeratorBot) CollectMessageAsEvidenceFromMessageModal(i *discordgo.InteractionCreate) {
	data := i.Interaction.ModalSubmitData()

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

// App command (where the target is a message), returns User reputation
func (bot *ModeratorBot) CheckUserReputation(i *discordgo.InteractionCreate) (reputation string, err error) {
	if i.Interaction.Member.User.ID == "" {
		return "", fmt.Errorf("user was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.CheckUserReputation,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation),
		},
	})

	return "", err
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

// UpdateModeratedUser updates moderated user status in the database.
// It is allowed to fail
func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Updates(u)

	if tx.RowsAffected != 1 {
		return fmt.Errorf("did not update one user row as expected, updated %v rows for user %s(%s)", tx.RowsAffected, u.UserName, u.UserID)
	}
	return nil
}
