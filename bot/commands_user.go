package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// TODO: I think all of these need to log events
func (bot *ModeratorBot) ModeratePositivelyFromUserContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User == nil {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.RespondToModeratePositivelyModalFromUserContext,
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

func (bot *ModeratorBot) ModerateNegativelyFromUserContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User == nil {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.RespondToModeratePositivelyModalFromUserContext,
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

// App command (where the target is a message), returns User reputation
func (bot *ModeratorBot) GetUserInfoFromUserContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: bot.userInfoIntegrationresponse(i),
	})
	if err != nil {
		log.Errorf("error responding to user info (message context), err: %v", err)
	}
}
