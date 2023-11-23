package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	globals "github.com/tyzbit/go-discord-modtools/globals"
)

// lowReputationValues returns a []discordgo.SelectMenuOption for low rep values
func lowReputationValues(sc ServerConfig) (options []discordgo.SelectMenuOption) {
	for i := -10; i <= 10; i++ {

		description := ""
		if sc.NotifyWhenReputationIsBelowSetting.Valid && int32(i) == sc.NotifyWhenReputationIsBelowSetting.Int32 {
			description = "Current value"
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprint(i),
			Description: description,
			Value:       fmt.Sprint(i),
		})
	}
	return options

}

// highReputationValues returns a []discordgo.SelectMenuOption for high rep values
func highReputationValues(sc ServerConfig) (options []discordgo.SelectMenuOption) {
	for i := -10; i <= 10; i++ {

		description := ""
		if sc.NotifyWhenReputationIsAboveSetting.Valid && int32(i) == sc.NotifyWhenReputationIsAboveSetting.Int32 {
			description = "Current value"
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprint(i),
			Description: description,
			Value:       fmt.Sprint(i),
		})
	}
	return options

}

// SettingsIntegrationResponse returns server settings in a *discordgo.InteractionResponseData
func (bot *ModeratorBot) SettingsIntegrationResponse(sc ServerConfig) *discordgo.InteractionResponseData {
	channel, _ := bot.DG.Channel(sc.EvidenceChannelSettingID)
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: getTagValue(sc, "NotifyWhenReputationIsBelowSetting", "pretty") + fmt.Sprintf(": %v", sc.NotifyWhenReputationIsBelowSetting.Int32),
						CustomID:    globals.NotifyWhenReputationIsBelowSetting,
						Options:     lowReputationValues(sc),
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: getTagValue(sc, "NotifyWhenReputationIsAboveSetting", "pretty") + fmt.Sprintf(": %v", sc.NotifyWhenReputationIsAboveSetting.Int32),
						CustomID:    globals.NotifyWhenReputationIsAboveSetting,
						Options:     highReputationValues(sc),
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder:  globals.EvidenceChannelSettingID + ": #" + channel.Name,
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
