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
	_, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannel.ID, &ms)
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
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{User: *m.Author}).First(&user)

	if user.User.ID == "" {
		return fmt.Errorf("unable to look up user %s(%s)", user.User.Username, &user.User.ID)
	}

	user.Reputation = user.Reputation + 1
	bot.UpdateModeratedUser(user)
	return nil
}

// App command (where the target is a message), copies message details to a
// configured channel then decreases the message author's reputation
func (bot *ModeratorBot) CollectMessageAsEvidenceThenDecreaseReputation(m *discordgo.Message) error {
	if m == nil {
		return fmt.Errorf("message was not provied")
	}

	// TODO: mirror message to other channel
	// TODO: decrease user rep
	return nil
}

// App command (where the target is a message), returns User reputation
func (bot *ModeratorBot) CheckUserReputationUsingMessage(m *discordgo.Message) (reputation int64, err error) {
	if m == nil {
		return 0, fmt.Errorf("message was not provied")
	}

	// TODO: look up user in DB
	// TODO: return embed with info
	return 0, nil
}

// App command (where the target is a user), returns User reputation
func (bot *ModeratorBot) CheckUserReputation(u *discordgo.User) (reputation int64, err error) {
	if u == nil {
		return 0, fmt.Errorf("message was not provied")
	}

	// TODO: look up user in DB
	// TODO: return embed with info
	return 0, nil
}

func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	// TODO: Update user in DB

	return nil
}
