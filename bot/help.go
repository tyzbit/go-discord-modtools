package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	// Slash command
	Help = "help"

	BotHelpText = `**Usage**
		In the Discord app, right click (Desktop) or long press (mobile) a message or user to see the available options.
	
	Configure the bot (Highly recommended to do so only people in a specific role can use the moderation commands):
	
	` + "`/settings`" + `
	
	Look up a user (more coming soon):
	
	` + "`/query`" + `
	
	Add a custom command (simply posts your desired block of text, Markdown formatting enabled)
	
	` + "`/addcommand`" + `
	
	Remove a custom command
	
	` + "`/deletecommand`" + `
	
	Get this help message:
	
	` + "`/help`"
)

// This produces the help text seen from the chat commant `/help`
func (bot *ModeratorBot) GetHelpFromChatCommandContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Moder8or Bot Help",
					Description: BotHelpText,
					Color:       Purple,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("error responding to help command "+Help+", err: %v", err)
	}
}
