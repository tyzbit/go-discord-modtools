package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// /moderate slash command, it needs a *discordgo.User at a minimum, either by
// direct reference or in relation to a *discordgo.Message
func (bot *ModeratorBot) Moderate(i *discordgo.InteractionCreate) error {
	if i.Interaction.Member.User == nil {
		return fmt.Errorf("user was not provided")
	} else if i.Message == nil {
		return fmt.Errorf("message was not provided")
	}

	// TODO: show embed with options
	return nil
}

// App command, copies message details to a configured channel
func (bot *ModeratorBot) CollectMessageAsEvidence(i *discordgo.InteractionCreate) error {
	if i.Message == nil {
		return fmt.Errorf("message was not provied")
	}

	sc := bot.getServerConfig(i.Message.GuildID)
	ms := discordgo.MessageSend{
		Content:    i.Message.Content,
		Embeds:     i.Message.Embeds,
		TTS:        i.Message.TTS,
		Components: i.Message.Components,
		//Files: m.Attachments,
		// AllowedMentions,
		Reference: i.Message.MessageReference,
		//File: ,
		// Embed: m.Embeds[],
	}
	_, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		return err
	}

	return nil
}

// App command (where the target is a message), copies message details to a
// configured channel then increases the message author's reputation
func (bot *ModeratorBot) CollectMessageAsEvidenceThenIncreaseReputation(i *discordgo.InteractionCreate) error {
	if i.Message == nil {
		return fmt.Errorf("message was not provied")
	}

	err := bot.CollectMessageAsEvidence(i)
	if err != nil {
		return err
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

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
func (bot *ModeratorBot) CollectMessageAsEvidenceThenDecreaseReputation(i *discordgo.InteractionCreate) error {
	if i.Message == nil {
		return fmt.Errorf("message was not provied")
	}

	err := bot.CollectMessageAsEvidence(i)
	if err != nil {
		return err
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

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
func (bot *ModeratorBot) CheckUserReputation(i *discordgo.InteractionCreate) (reputation int64, err error) {
	if i.Interaction.Member.User.ID == "" {
		return 0, fmt.Errorf("message was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	if user.UserID == "" {
		return 0, fmt.Errorf("unable to look up user %s(%s)", user.UserName, user.UserID)
	}

	// TODO: I'm testing modals here for some reason but this response
	// should probably just be an ephemeral message.
	err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "12345",
			Title:    "Modal test",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "12345",
						Label:       "Reputation",
						Style:       discordgo.TextInputShort,
						Placeholder: "this user is amazing",
						Required:    true,
						MinLength:   1,
						MaxLength:   4,
					},
				},
			}},
		},
	})

	return user.Reputation, nil
}

func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Updates(u)

	if tx.RowsAffected != 1 {
		return fmt.Errorf("did not update one user row as expected, updated %v rows for user %s(%s)", tx.RowsAffected, u.UserName, u.UserID)
	}
	return nil
}
