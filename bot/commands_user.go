package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// TODO: I think all of these need to log events
func (bot *ModeratorBot) ModerateFromUserContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User == nil {
		log.Warn("no user nor message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ShowModerationModalFromUserContext,
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
func (bot *ModeratorBot) GetUserInfoFromUserContext(i *discordgo.InteractionCreate) (reputation string, err error) {
	if i.Interaction.Member.User.ID == "" {
		return "", fmt.Errorf("user was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Message.Author.ID}).First(&user)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.GetUserInfoFromUserContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation),
		},
	})

	return "", err
}
