package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Responds to an InteractionCreate with a permission denied message
func (bot *ModeratorBot) RespondWithPermissionDenied(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: bot.permissionsErrorDisplayedToTheUser(),
	})
	if err != nil {
		log.Warn("error responding to interaction message interaction: %w", err)
	}
}

// SettingsIntegrationResponse returns server settings in a *discordgo.InteractionResponseData
func (bot *ModeratorBot) SettingsIntegrationResponse(cfg GuildConfig) *discordgo.InteractionResponseData {
	channel, _ := bot.DG.Channel(cfg.EvidenceChannelSettingID)
	var evidenceChannelName, moderatorRoleName string
	if channel == nil {
		evidenceChannelName = "not set"
	} else {
		evidenceChannelName = "#" + channel.Name
	}
	moderatorRole, _ := bot.DG.State.Role(cfg.ID, cfg.ModeratorRoleSettingID)
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
						Placeholder:  EvidenceChannelSettingID + ": " + evidenceChannelName,
						MenuType:     discordgo.ChannelSelectMenu,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						CustomID:     EvidenceChannelSettingID,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: ModeratorRoleSettingID + ": " + moderatorRoleName,
						MenuType:    discordgo.RoleSelectMenu,
						CustomID:    ModeratorRoleSettingID,
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

	user := bot.GetModeratedUser(i.GuildID, i.Interaction.Member.User.ID)

	return &discordgo.InteractionResponseData{
		CustomID: GetUserInfoFromUserContext,
		Flags:    discordgo.MessageFlagsEphemeral,
		Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, *user.Reputation),
	}
}

// Simple wrapper to display an embed to the user with an error (ephemeral)
func (bot *ModeratorBot) generalErrorDisplayedToTheUser(reason string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "There was an issue",
				Description: reason,
				Color:       DarkRed,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}

// Simple wrapper to display an embed to the user with an error (ephemeral)
func (bot *ModeratorBot) generalInfoDisplayedToTheUser(reason string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Notice",
				Description: reason,
				Color:       FrenchGray,
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
				Color: DarkRed,
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}

// Simple wrapper to display an embed to the user with an error (ephemeral)
func (bot *ModeratorBot) settingsErrorDisplayedToTheUser() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Please configure the bot first",
				Color: DarkRed,
				Description: "Please use /" + Settings + " to set " +
					"the mod role and evidence channel",
			},
		},
		Flags: discordgo.MessageFlagsEphemeral,
	}
}
