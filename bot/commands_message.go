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
		log.Warn("message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.RespondToModerationModalFromMessageContext,
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
							MenuType:    discordgo.UserSelectMenu,
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
						discordgo.SelectMenu{
							Placeholder:  "Channel",
							MenuType:     discordgo.ChannelSelectMenu,
							ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
							CustomID:     message.ChannelID,
							Options: []discordgo.SelectMenuOption{{
								Label: "Channel",
								Value: "Saved",
							}},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder: "Message",
							CustomID:    message.ID,
							Options: []discordgo.SelectMenuOption{{
								Label: "Message",
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

func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	return
}

func (bot *ModeratorBot) SaveEvidenceFromMessageContext(i *discordgo.InteractionCreate) {
	sc := bot.getServerConfig(i.GuildID)

	// TODO: more information
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Evidence saved to <#" + sc.EvidenceChannelSettingID + ">",
		},
	})
}
