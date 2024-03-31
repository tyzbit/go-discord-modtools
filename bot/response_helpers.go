package bot

import (
	"github.com/bwmarrin/discordgo"
)

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
