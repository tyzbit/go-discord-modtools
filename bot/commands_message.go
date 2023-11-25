package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Moderation modal menu
func (bot *ModeratorBot) DocumentBehaviorFromMessageContext(i *discordgo.InteractionCreate) {
	data := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if data.ID == "" {
		reason := "No message was provided"
		log.Warn(reason)
		_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bot.generalErrorDisplayedToTheUser(reason)})
		return
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fmt.Sprintf("Log user behavior for %s (ID: %v)", data.Author.Username, data.Author.ID),
					Description: "Document user behavior - good, bad, or noteworthy",
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
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
