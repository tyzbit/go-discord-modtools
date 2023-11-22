package bot

import (
	"fmt"

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
					discordgo.Button{
						Label: getTagValue(sc, "NotifyWhenReputationIsBelowSetting", "pretty") +
							fmt.Sprintf(", current: %v", sc.NotifyWhenReputationIsBelowSetting.Int32),
						Style:    discordgo.PrimaryButton,
						CustomID: globals.NotifyWhenReputationIsBelowSetting},
					discordgo.Button{
						Label: getTagValue(sc, "NotifyWhenReputationIsAboveSetting", "pretty") +
							fmt.Sprintf(", current: %v", sc.NotifyWhenReputationIsAboveSetting.Int32),
						Style:    discordgo.PrimaryButton,
						CustomID: globals.NotifyWhenReputationIsAboveSetting},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder:  getTagValue(sc, "EvidenceChannelSettingID", "pretty"),
						MenuType:     discordgo.ChannelSelectMenu,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						CustomID:     globals.EvidenceChannelSettingID,
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
				Color: globals.Purple,
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
				Color: globals.Purple,
			},
		},
	}
}
