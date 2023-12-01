package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	globals "github.com/tyzbit/go-discord-modtools/globals"
)

// SettingsIntegrationResponse returns server settings in a *discordgo.InteractionResponseData
func (bot *ModeratorBot) SettingsIntegrationResponse(sc ServerConfig) *discordgo.InteractionResponseData {
	channel, _ := bot.DG.Channel(sc.EvidenceChannelSettingID)
	var evidenceChannelName, moderatorRoleName string
	if channel == nil {
		evidenceChannelName = "not set"
	} else {
		evidenceChannelName = "#" + channel.Name
	}
	moderatorRole, _ := bot.DG.State.Role(sc.DiscordId, sc.ModeratorRoleSettingID)
	if moderatorRole == nil {
		moderatorRoleName = "not set"
	} else {
		moderatorRoleName = moderatorRole.Name
	}
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder:  globals.EvidenceChannelSettingID + ": " + evidenceChannelName,
						MenuType:     discordgo.ChannelSelectMenu,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						CustomID:     globals.EvidenceChannelSettingID,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: globals.ModeratorRoleSettingID + ": " + moderatorRoleName,
						MenuType:    discordgo.RoleSelectMenu,
						CustomID:    globals.ModeratorRoleSettingID,
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

// Simple wrapper to display an embed to the user with an error (ephemeral)
func (bot *ModeratorBot) permissionsErrorDisplayedToTheUser() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "You do not have permission to do that",
				Color: globals.DarkRed,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}
