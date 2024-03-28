package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Starts a poll the user has customized
func (bot *ModeratorBot) StartPoll(i *discordgo.InteractionCreate, msd discordgo.ModalSubmitInteractionData) {
	pollTitle := msd.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	var options []string
	// Intentionally starting at the __second__ component in the ModalSubmitInteractionData because the first is
	// the title of the poll
	for i := 1; i <= len(msd.Components)-1; i++ {
		options = append(options, msd.Components[i].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	}

	// Create fields for every option
	components := []discordgo.MessageComponent{}
	for i, option := range options {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{discordgo.Button{
				CustomID: fmt.Sprintf("%s%v", PollOptionPrefix, i),
				Label:    option,
				Style:    discordgo.PrimaryButton,
			}},
		})
	}

	// Add am "End poll" button
	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{discordgo.Button{
			CustomID: EndPoll,
			Style:    discordgo.DangerButton,
			Label:    EndPoll,
		}},
	})

	// Send message in channel instead
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       pollTitle,
					Description: fmt.Sprintf("Created by <@%s>", i.Interaction.Member.User.ID),
					Color:       Purple,
				},
			},
			Flags:      discordgo.MessageFlagsEphemeral, // TODO: Remove after testing
			Components: components,
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}
