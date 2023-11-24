package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

func (bot *ModeratorBot) RespondToModerationModalFromMessageContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Moderation action saved!",
					Description: globals.ModerationSuccessful,
					Color:       globals.Purple,
				},
			},
		},
	})
}

func (bot *ModeratorBot) RespondToModerationModalFromUserContext(i *discordgo.InteractionCreate) {
	// TODO: check status and change message base on status
	// Drop a message in the evidence channel
	bot.SaveEvidenceFromModalSubmission(i)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Moderation action saved!",
					Description: globals.ModerationSuccessful,
					Color:       globals.Purple,
				},
			},
		},
	})
}

// Modal command, copies message details to a configured channel
func (bot *ModeratorBot) SaveEvidenceFromModalSubmission(i *discordgo.InteractionCreate) {
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
