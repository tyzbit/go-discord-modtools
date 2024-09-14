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
	user := i.Interaction.ApplicationCommandData().Resolved.Users[i.ApplicationCommandData().TargetID]
	var err error
	var cfg GuildConfig
	bot.DB.Where(&GuildConfig{ID: i.GuildID}).First(&cfg)
	if cfg.EvidenceChannelSettingID == "" {
		err = bot.DG.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: bot.settingsErrorDisplayedToTheUser(),
			})
	} else {
		err = bot.DG.InteractionRespond(i.Interaction,
			bot.GenerateEvidenceReportFromUser(i, user))
	}
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
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
