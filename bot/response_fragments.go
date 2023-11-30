package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	globals "github.com/tyzbit/go-discord-modtools/globals"
)

// lowReputationValues returns a []discordgo.SelectMenuOption for the configurable
// value to notify when a reputation drops to
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

// highReputationValues returns a []discordgo.SelectMenuOption for the configurable
// value to notify when a reputation rises to
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

// User information and stats produced for the /query command and
// "Get info" when right clicking users
func (bot *ModeratorBot) userInfoIntegrationresponse(i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	if i.Interaction.Member.User.ID == "" {
		log.Warn("user was not provided")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	return &discordgo.InteractionResponseData{
		CustomID: globals.GetUserInfoFromUserContext,
		Flags:    discordgo.MessageFlagsEphemeral,
		Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation.Int64),
	}
}

// Simple wrapper to display an embed to the user with an error (ephemeral)
func (bot *ModeratorBot) generalErrorDisplayedToTheUser(reason string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "There was an issue",
				Description: reason,
				Color:       globals.DarkRed,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}
