package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Moderation modal menu
func (bot *ModeratorBot) SaveEvidenceFromMessageContext(i *discordgo.InteractionCreate) {
	return
}

func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User.ID == "" {
		log.Warn("user was not provided")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	// TODO: Add more user information
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.GetUserInfoFromUserContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation),
		},
	})
}

func (bot *ModeratorBot) IncreaseReputationFromMessageContext(i *discordgo.InteractionCreate) {
	// Currently, you can only use TextInput in modal action rows builders.
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.IncreaseReputationFromMessageContext,
			Title:    "Additional details",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.ReasonOption,
							Label:       "Reason",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Reason for moderation (optional, shown in configured evidence channel)",
							Required:    false,
							MaxLength:   300,
						},
					},
				},
			},
		},
	})
}

func (bot *ModeratorBot) DecreaseReputationFromMessageContext(i *discordgo.InteractionCreate) {
	// Currently, you can only use TextInput in modal action rows builders.
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.DecreaseReputationFromMessageContext,
			Title:    "Additional details",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.ReasonOption,
							Label:       "Reason",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Reason for moderation (optional, shown in configured evidence channel)",
							Required:    false,
							MaxLength:   300,
						},
					},
				},
			},
		},
	})
}
