package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Called from the App menu, this displays an embed for the moderator to
// choose to change the reputation of the posting user
// and (PLANNED) produces output in the evidence channel with information about
// the message, user and moderation actions taken
func (bot *ModeratorBot) DocumentBehaviorFromMessageContext(i *discordgo.InteractionCreate) {
	message := *i.Interaction.ApplicationCommandData().Resolved.
		Messages[i.ApplicationCommandData().TargetID]
	// This check might be redundant - we may never get here without message
	// ApplicationCommandData (unless we call this mistakenly from another context)
	if message.ID == "" {
		reason := "No message was provided"
		log.Warn(reason)
		err := bot.DG.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: bot.generalErrorDisplayedToTheUser(reason),
			},
		)
		if err != nil {
			log.Warn("error responding to interaction: %w", err)
		}
		return
	}

	var err error
	sc := bot.getServerConfig(i.GuildID)
	if sc.EvidenceChannelSettingID == "" {
		err = bot.DG.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: bot.settingsErrorDisplayedToTheUser(),
			})
	} else {
		err = bot.DG.InteractionRespond(i.Interaction,
			bot.GenerateEvidenceReportFromMessage(i, &message))
	}
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Produces user info such as reputation and (PLANNED) stats
func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bot.userInfoIntegrationresponse(i),
		})
	if err != nil {
		log.Errorf("error responding to user info (message context), err: %v", err)
	}
}

// Returns a discordgo.InteractionResponse with an evidence report based on a message provided
func (bot *ModeratorBot) GenerateEvidenceReportFromMessage(i *discordgo.InteractionCreate, message *discordgo.Message) (resp *discordgo.InteractionResponse) {
	user := bot.GetModeratedUser(i.GuildID, message.Author.ID)
	var fields []*discordgo.MessageEmbedField
	var messageType discordgo.InteractionResponseType
	var authorID string
	if len(message.Embeds) > 0 {
		fields = message.Embeds[0].Fields
		messageType = discordgo.InteractionResponseUpdateMessage
		authorID = getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[1].Value)
		for idx, field := range fields {
			if field.Name == CurrentReputation {
				user := bot.GetModeratedUser(i.GuildID, authorID)
				fields[idx].Value = fmt.Sprintf("%v", user.Reputation.Int64)
			}
		}
	} else {
		messageType = discordgo.InteractionResponseChannelMessageWithSource
		authorID = message.Author.ID
		fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Author of message",
				Value:  fmt.Sprintf("<@%s>", user.UserID),
				Inline: true,
			},
			{
				Name:   PreviousReputation,
				Value:  fmt.Sprintf("%v", user.Reputation.Int64),
				Inline: true,
			},
			{
				Name:   CurrentReputation,
				Value:  fmt.Sprintf("%v", user.Reputation.Int64),
				Inline: true,
			},
			{
				Name:   "Link to original message",
				Value:  fmt.Sprintf(MessageURLTemplate, i.Interaction.GuildID, message.ChannelID, message.ID),
				Inline: true,
			},
			{
				Name:  OriginalMessageContent,
				Value: message.Content,
			},
			{
				Name:   "Original message timestamp",
				Value:  fmt.Sprintf("%s (<t:%v:R>)", message.Timestamp.Format(time.RFC1123Z), message.Timestamp.Unix()),
				Inline: false,
			},
		}

		if len(message.Attachments) > 0 {
			attachmentList := ""
			for _, attachment := range message.Attachments {
				attachmentList = attachmentList + fmt.Sprintf("[%s](%v)\n", attachment.Filename, attachment.URL)
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf(Attachments+" (%v)", len(message.Attachments)),
				Value: attachmentList,
			})
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "Collected by",
			Value: fmt.Sprintf(`<@%s>
									%s (<t:%v:R>)`,
				i.Interaction.Member.User.ID,
				time.Now().UTC().Format(time.RFC1123Z),
				time.Now().Unix()),
		})
	}

	return &discordgo.InteractionResponse{
		Type: messageType,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Evidence Report",
					Description: fmt.Sprintf("Document user behavior for for <@%v> - good, bad, or noteworthy", authorID),
					Color:       Purple,
					Fields:      fields,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: DecreaseUserReputation,
						Label:    DecreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: IncreaseUserReputation,
						Label:    IncreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: ShowEvidenceCollectionModal,
						Label:    ShowEvidenceCollectionModal,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: SubmitReport,
						Label:    SubmitReport,
						Style:    discordgo.PrimaryButton,
					},
				},
			}},
		},
	}
}
