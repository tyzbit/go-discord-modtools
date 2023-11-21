package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// /moderate slash command, it needs a *discordgo.User at a minimum, either by
// direct reference or in relation to a *discordgo.Message
func (bot *ModeratorBot) Moderate(u *discordgo.User, m *discordgo.Message) error {
	if u == nil {
		return fmt.Errorf("user was not provided")
	} else if m == nil {
		return fmt.Errorf("message was not provided")
	}

	// TODO: show embed with options
	return nil
}

// App command, copies message details to a configured channel
func (bot *ModeratorBot) CollectMessageAsEvidence(m *discordgo.Message) error {
	if m == nil {
		return fmt.Errorf("message was not provied")
	}

	sc := bot.getServerConfig(m.GuildID)
	ms := discordgo.MessageSend{
		Content:    m.Content,
		Embeds:     m.Embeds,
		TTS:        m.TTS,
		Components: m.Components,
		//Files: m.Attachments,
		// AllowedMentions,
		Reference: m.MessageReference,
		//File: ,
		// Embed: m.Embeds[],
	}
	_, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelID, &ms)
	if err != nil {
		return err
	}

	return nil
}

// App command (where the target is a message), copies message details to a
// configured channel then increases the message author's reputation
func (bot *ModeratorBot) CollectMessageAsEvidenceThenIncreaseReputation(m *discordgo.Message) error {
	if m == nil {
		return fmt.Errorf("message was not provied")
	}

	err := bot.CollectMessageAsEvidence(m)
	if err != nil {
		return err
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: m.Author.ID}).First(&user)

	if user.UserID == "" {
		return fmt.Errorf("unable to look up user %s(%s)", user.UserName, user.UserID)
	}

	user.Reputation = user.Reputation + 1
	err = bot.UpdateModeratedUser(user)
	if err != nil {
		return err
	}
	return nil
}

// App command (where the target is a message), copies message details to a
// configured channel then decreases the message author's reputation
func (bot *ModeratorBot) CollectMessageAsEvidenceThenDecreaseReputation(m *discordgo.Message) error {
	if m == nil {
		return fmt.Errorf("message was not provied")
	}

	err := bot.CollectMessageAsEvidence(m)
	if err != nil {
		return err
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: m.Author.ID}).First(&user)

	if user.UserID == "" {
		return fmt.Errorf("unable to look up user %s(%s)", user.UserName, user.UserID)
	}

	user.Reputation = user.Reputation - 1
	err = bot.UpdateModeratedUser(user)
	if err != nil {
		return err
	}
	return nil
}

// App command (where the target is a message), returns User reputation
func (bot *ModeratorBot) CheckUserReputationUsingMessage(m *discordgo.Message) (reputation int64, err error) {
	if m == nil {
		return 0, fmt.Errorf("message was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: m.Author.ID}).First(&user)

	if user.UserID == "" {
		return 0, fmt.Errorf("unable to look up user %s(%s)", user.UserName, user.UserID)
	}

	return user.Reputation, nil
}

// App command (where the target is a user), returns User reputation
func (bot *ModeratorBot) CheckUserReputation(u *discordgo.User) (reputation int64, err error) {
	if u == nil {
		return 0, fmt.Errorf("message was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.ID}).First(&user)

	if user.UserID == "" {
		return 0, fmt.Errorf("unable to look up user %s(%s)", user.UserName, user.UserID)
	}

	return user.Reputation, nil
}

func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Updates(u)

	if tx.RowsAffected != 1 {
		return fmt.Errorf("did not update one user row as expected, updated %v rows for user %s(%s)", tx.RowsAffected, u.UserName, u.UserID)
	}
	return nil
}
