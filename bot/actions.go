package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// TODO: I think all of these need to log events

// /moderate slash command, it needs a *discordgo.User at a minimum, either by
// direct reference or in relation to a *discordgo.Message
func (bot *ModeratorBot) Moderate(i *discordgo.InteractionCreate) error {
	if i.Interaction.Member.User == nil {
		return fmt.Errorf("user was not provided")
	} else if i.Interaction.Message == nil {
		return fmt.Errorf("message was not provided")
	}

	mcd := i.MessageComponentData()
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.ModerateModal,
			Title:    "Moderate " + mcd.Values[0],
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    globals.ModerateModalReason,
						Label:       "Reason",
						Style:       discordgo.TextInputShort,
						Placeholder: "Why are you moderating this user or content?",
						Required:    true,
						MinLength:   1,
						MaxLength:   500,
					},
				},
			}},
		},
	})
	return nil
}

// App command, copies message details to a configured channel
func (bot *ModeratorBot) CollectMessageAsEvidence(i *discordgo.InteractionCreate) error {
	message := i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if message == nil {
		return fmt.Errorf("message was not provied")
	}

	// sc := bot.getServerConfig(i.Interaction.GuildID)
	// ms := discordgo.MessageSend{
	// 	Content: message.Content,
	// 	Embeds: []*discordgo.MessageEmbed{{
	// 		Fields: []*discordgo.MessageEmbedField{
	// 			{
	// 				Name:   "Username",
	// 				Value:  message.Author.Username,
	// 				Inline: true,
	// 			},
	// 			{
	// 				Name:   "Channel",
	// 				Value:  fmt.Sprintf("<#" + message.ChannelID + ">"),
	// 				Inline: true,
	// 			},
	// 			{
	// 				Name:   "Timestamp",
	// 				Value:  message.Timestamp.Format(time.RFC1123Z),
	// 				Inline: false,
	// 			},
	// 		},
	// 	}},
	// 	TTS:        message.TTS,
	// 	Components: message.Components,
	// 	//Files: m.Attachments,
	// 	// AllowedMentions,
	// 	Reference: message.MessageReference,
	// 	//File: ,
	// 	// Embed: m.Embeds[],
	// }
	// _, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	// if err != nil {
	// 	return err
	// }

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message.Content,
			Embeds: []*discordgo.MessageEmbed{{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Username",
						Value:  message.Author.Username,
						Inline: true,
					},
					{
						Name:   "Channel",
						Value:  fmt.Sprintf("<#" + message.ChannelID + ">"),
						Inline: true,
					},
					{
						Name:   "Timestamp",
						Value:  message.Timestamp.Format(time.RFC1123Z),
						Inline: false,
					},
				},
			}},
			TTS:        message.TTS,
			Components: message.Components,
			//Files: m.Attachments,
			// AllowedMentions,
			//File: ,
			// Embed: m.Embeds[],

		},
	})

	return nil
}

// App command (where the target is a message), copies message details to a
// configured channel then increases the message author's reputation
func (bot *ModeratorBot) CollectMessageAsEvidenceThenIncreaseReputation(i *discordgo.InteractionCreate) error {
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
func (bot *ModeratorBot) CheckUserReputation(i *discordgo.InteractionCreate) (reputation string, err error) {
	if i.Interaction.Member.User.ID == "" {
		return "", fmt.Errorf("user was not provied")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.CheckUserReputation,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation),
		},
	})

	return "", err
}

func (bot *ModeratorBot) ShowLowReputationModal(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.NotifyWhenReputationIsBelowSetting,
			Title:    "What should the new value be?",
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
}

func (bot *ModeratorBot) ShowHighReputationModal(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.NotifyWhenReputationIsAboveSetting,
			Title:    "What should the new value be?",
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
}

// func (bot *ModeratorBot) SetValueUsingModal(i *discordgo.InteractionCreate) {
// 	guild, err := bot.DG.Guild(i.Interaction.GuildID)
// 	if err != nil {
// 		log.Errorf("unable to look up guild ID %s", i.Interaction.GuildID)
// 		return
// 	}

// 	// The first of two events because after choosing the option, the user
// 	// inputs a value in a modal.
// 	bot.createInteractionEvent(InteractionEvent{
// 		UserID:        i.Member.User.ID,
// 		Username:      i.Member.User.Username,
// 		InteractionId: i.Message.ID,
// 		ChannelId:     i.Message.ChannelID,
// 		ServerID:      i.Interaction.GuildID,
// 		ServerName:    guild.Name,
// 	})

// 	mcd := i.MessageComponentData()
// 	if mcd.CustomID == globals.ShowLowReputationModal {
// 		bot.respondToSettingsChoice(i, "notify_when_reputation_is_below_setting", mcd.Values[0])
// 	} else if mcd.CustomID == globals.ShowHighReputationModal {
// 		bot.respondToSettingsChoice(i, "notify_when_reputation_is_above_setting", mcd.Values[0])
// 	}
// }

// UpdateModeratedUser updates moderated user status in the database.
// It is allowed to fail
func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Updates(u)

	if tx.RowsAffected != 1 {
		return fmt.Errorf("did not update one user row as expected, updated %v rows for user %s(%s)", tx.RowsAffected, u.UserName, u.UserID)
	}
	return nil
}
