package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Moderation modal menu
func (bot *ModeratorBot) SaveEvidenceFromMessageContext(i *discordgo.InteractionCreate) {
	return
}

func (bot *ModeratorBot) ModeratePositivelyFromMessageContext(i *discordgo.InteractionCreate) {
	// Currently, you can only use TextInput in modal action rows builders.
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ModeratePositivelyFromMessageContext,
			Title:    "Additional details (Moderating positively)",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.ReasonOption,
							Label:       "Reason",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Reason for moderation (optional, posted to evidence channel)",
							Required:    false,
							MaxLength:   300,
						},
					},
				},
			},
		},
	})
}

func (bot *ModeratorBot) ModerateNegativelyFromMessageContext(i *discordgo.InteractionCreate) {
	// Currently, you can only use TextInput in modal action rows builders.
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ModerateNegativelyFromMessageContext,
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

func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: bot.userInfoIntegrationresponse(i),
	})
	if err != nil {
		log.Errorf("error responding to user info (message context), err: %v", err)
	}
}
