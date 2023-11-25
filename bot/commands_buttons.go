package bot

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// RespondToSettingsChoice updates a server setting according to the
// column name (setting) and the value
func (bot *ModeratorBot) RespondToSettingsChoice(i *discordgo.InteractionCreate,
	setting string, value interface{}) {
	guild, err := bot.DG.Guild(i.Interaction.GuildID)
	if err != nil {
		log.Errorf("unable to look up guild ID %s", i.Interaction.GuildID)
		return
	}

	sc, ok := bot.updateServerSetting(i.Interaction.GuildID, setting, value)
	var interactionErr error

	bot.createInteractionEvent(InteractionEvent{
		UserID:        i.Member.User.ID,
		Username:      i.Member.User.Username,
		InteractionId: i.Message.ID,
		ChannelId:     i.Message.ChannelID,
		ServerID:      i.Interaction.GuildID,
		ServerName:    guild.Name,
	})

	if !ok {
		reason := "Unable to save settings"
		interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: bot.generalErrorDisplayedToTheUser(reason),
		})
	} else {
		interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: bot.SettingsIntegrationResponse(sc),
		})
	}

	if interactionErr != nil {
		log.Errorf("error responding to settings interaction, err: %v", interactionErr)
	}
}

func (bot *ModeratorBot) ChangeUserReputation(i *discordgo.InteractionCreate, increase bool) {
	pattern := regexp.MustCompile(`<@(\d+)>`)

	userID := ""
	match := pattern.FindStringSubmatch(i.Interaction.Message.Embeds[0].Fields[1].Value)
	if len(match) > 1 {
		userID = match[1]
	} else {
		log.Warnf("unable to get user ID from message ID %s", i.Interaction.Message.ID)
		return
	}

	user := ModeratedUser{}
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: userID}).First(&user)

	if tx.RowsAffected > 1 {
		log.Errorf("unexpected number of rows affected getting user reputation: %v", tx.RowsAffected)
		return
	} else if tx.RowsAffected == 0 {
		user.UserID = userID
	}

	if increase {
		user.Reputation = user.Reputation + 1
	} else {
		user.Reputation = user.Reputation - 1
	}

	err := bot.UpdateModeratedUser(user)
	if err != nil {
		log.Warn("unable to update user moderation record, err: %w", err)
	}
}
