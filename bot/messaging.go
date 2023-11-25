package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// This whole file can probably be deleted (artifacts from copying from the other project)

// sendModerateResponse sends the message with a result from Moder8or
func (bot *ModeratorBot) sendModerateResponse(userMessage *discordgo.Message, messagesToSend *discordgo.MessageSend) error {
	username := ""
	user, err := bot.DG.User(userMessage.Member.User.ID)
	if err != nil {
		log.Errorf("unable to look up user with ID %v, err: %v", userMessage.Member.User.ID, err)
		username = "unknown"
	} else {
		username = user.Username
	}

	var guild *discordgo.Guild
	if userMessage.GuildID != "" {
		var gErr error
		// Do a lookup for the full guild object
		guild, gErr = bot.DG.Guild(userMessage.GuildID)
		if gErr != nil {
			return gErr
		}
		bot.createMessageEvent(MessageEvent{
			AuthorId:       user.ID,
			AuthorUsername: user.Username,
			MessageId:      userMessage.ID,
			ChannelId:      userMessage.ChannelID,
			ServerID:       userMessage.GuildID,
		})
		log.Debugf("sending moderation message response in %s(%s), calling user: %s(%s)",
			guild.Name, guild.ID, username, userMessage.Member.User.ID)
	}

	botMessage, err := bot.DG.ChannelMessageSendComplex(userMessage.ChannelID, messagesToSend)
	// For some reason, this message is absent a Guild ID, so we copy from the previous message
	if guild.ID != "" {
		botMessage.GuildID = guild.ID
	}

	if err != nil {
		log.Errorf("problem sending message: %v", err)
		return err
	}

	return nil
}

// sendModerateResponse sends the message with a result from Moder8or
func (bot *ModeratorBot) sendModerateCommandResponse(i *discordgo.Interaction, message *discordgo.MessageSend) error {
	username := ""
	var user *discordgo.User
	var err error
	if i.User != nil {
		user, err = bot.DG.User(i.User.ID)
	} else {
		user, err = bot.DG.User(i.Member.User.ID)
	}
	if err != nil {
		log.Errorf("unable to look up user with ID %v, err: %v", i.User.ID, err)
		username = "unknown"
	} else {
		username = user.Username
	}

	if i.GuildID != "" {
		// Do a lookup for the full guild object
		guild, gErr := bot.DG.Guild(i.GuildID)
		if gErr != nil {
			return gErr
		}
		bot.createMessageEvent(MessageEvent{
			AuthorId:       user.ID,
			AuthorUsername: user.Username,
			MessageId:      i.ID,
			ChannelId:      i.ChannelID,
			ServerID:       i.GuildID,
		})
		log.Debugf("sending moderation message response in %s(%s), calling user: %s(%s)",
			guild.Name, guild.ID, username, user.ID)
	}

	interactionMessage, err := bot.DG.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds:     &message.Embeds,
		Components: &message.Components,
	})

	if err != nil {
		return err
	}

	// For some reason, this message is absent a Guild ID, so we copy from the previous message
	if i.GuildID != "" {
		interactionMessage.GuildID = i.GuildID
	}

	// We don't remove the reply button because one, the message is visible only
	// to the calling user, so the space it takes up shouldn't matter (they
	// can dismiss the message entirely as well). Second, it doesn't seem it's
	// possible to edit that kind of message ¯\_(ツ)_/¯
	return nil
}
