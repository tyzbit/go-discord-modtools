package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// TODO: I think all of these need to log events

// Very similar to GenerateEvidenceReportFromMessageContext, but this is
// called when the target is a user and not a message, therefore
// this will be implicitly without any message reference
func (bot *ModeratorBot) DocumentBehaviorFromUserContext(i *discordgo.InteractionCreate) {
	// if i.Interaction.Member.User == nil {
	// 	reason := "No user was provided"
	// 	log.Warn(reason)
	// 	err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: bot.generalErrorDisplayedToTheUser(reason)})
	// 	return
	// }

	// data := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	// err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Embeds: []*discordgo.MessageEmbed{
	// 			{
	// 				Title:       fmt.Sprintf("Log user behavior for %s (ID: %v)", data.Author.Username, data.Author.ID),
	// 				Description: "Document user behavior, good bad or noteworthy",
	// 			},
	// 		},
	// 		Flags: discordgo.MessageFlagsEphemeral,
	// 	},
	// })
}

// Produces user info such as reputation and (PLANNED) stats
func (bot *ModeratorBot) GetUserInfoFromUserContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: bot.userInfoIntegrationresponse(i),
	})
	if err != nil {
		log.Errorf("error responding to user info (message context), err: %v", err)
	}
}
