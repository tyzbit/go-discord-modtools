package bot

import (
	"github.com/bwmarrin/discordgo"
	globals "github.com/tyzbit/go-discord-modtools/globals"
)

// SettingsIntegrationResponse returns server settings in a *discordgo.InteractionResponseData
func (bot *ModeratorBot) SettingsIntegrationResponse(sc ServerConfig) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						Label:    getTagValue(sc, "NotifyWhenReputationIsBelowSetting", "pretty"),
						Style:    discordgo.TextInputShort,
						CustomID: globals.NotifyWhenReputationIsBelowSetting,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						Label:    getTagValue(sc, "NotifyWhenReputationIsAboveSetting", "pretty"),
						Style:    discordgo.TextInputShort,
						CustomID: globals.NotifyWhenReputationIsAboveSetting,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: getTagValue(sc, "EvidenceChannelSetting", "pretty"),
						MenuType:    discordgo.ChannelSelectMenu,
						CustomID:    globals.EvidenceChannelSetting,
					},
				},
			},
		},
	}
}

// settingsFailureIntegrationResponse returns a *discordgo.InteractionResponseData
// stating that a failure to update settings has occured
func (bot *ModeratorBot) settingsFailureIntegrationResponse() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Unable to update setting",
				Color: globals.FrenchGray,
			},
		},
	}
}

// settingsFailureIntegrationResponse returns a *discordgo.InteractionResponseData
// stating that a failure to update settings has occured
func (bot *ModeratorBot) settingsDMFailureIntegrationResponse() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "The bot does not have any per-user settings",
				Color: globals.FrenchGray,
			},
		},
	}
}
