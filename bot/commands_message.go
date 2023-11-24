package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Moderation modal menu
func (bot *ModeratorBot) ModerateFromMessageContext(i *discordgo.InteractionCreate) {
	message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if message.ID == "" {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ShowModerationModalFromMessageContext,
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
	// 		CustomID: globals.ModerateUserFromModalContext,
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

// App command, copies message details to a configured channel
// TODO: fill out function
func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {

}

// App command, copies message details to a
// TODO: fill out function
func (bot *ModeratorBot) SaveEvidenceFromModalSubmissionFromMessageContext(i *discordgo.InteractionCreate) {

}
