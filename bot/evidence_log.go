package bot

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	// User commands
	GetUserInfoFromUserContext      = "Check info"
	DocumentBehaviorFromUserContext = "Save evidence"

	// Message commands
	GetUserInfoFromMessageContext      = "Get user info"
	DocumentBehaviorFromMessageContext = "Save as evidence"

	// Settings (the names affect the column names in the DB)
	EvidenceChannelSettingID = "Evidence channel"
	ModeratorRoleSettingID   = "Moderator role"

	// Moderation buttons
	IncreaseUserReputation      = "⬆️ Reputation"
	DecreaseUserReputation      = "⬇️ Reputation"
	ShowEvidenceCollectionModal = "Add notes"
	SubmitReport                = "Submit report"

	// Modals
	SaveEvidenceNotes = "Save evidence notes"

	// Modal options
	EvidenceNotes = "Evidence notes"

	// Text fragments
	CurrentReputation      = "Current reputation"
	PreviousReputation     = "Previous reputation"
	ModerationSuccessful   = "Moderation action saved"
	ModerationUnSuccessful = "There was a problem saving moderation action"
	OriginalMessageContent = "Content of original message"
	Attachments            = "Attachments"
	Notes                  = "Notes"
)

// Called from the App menu, this displays an embed for the moderator to
// choose to change the reputation of the posting user
// and (PLANNED) produces output in the evidence channel with information about
// the message, user and moderation actions taken
func (bot *ModeratorBot) DocumentBehaviorFromButtonContext(i *discordgo.InteractionCreate) {
	message := *i.Interaction.Message
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

	err := bot.DG.InteractionRespond(i.Interaction,
		bot.GenerateEvidenceReportFromMessage(i, &message))
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// ShowEvidenceCollectionModal shows the user a modal to fill in information
// about selected evidence
func (bot *ModeratorBot) ShowEvidenceCollectionModal(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: SaveEvidenceNotes,
			Title:    "Add notes to this report",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    EvidenceNotes,
						Label:       "Notes",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "Add notes such as reasoning or context",
						Required:    true,
						MinLength:   3,
						MaxLength:   MaxMessageContentLength,
					},
				},
			}},
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Submits evidence to the configured channel (including notes)
func (bot *ModeratorBot) SubmitReport(i *discordgo.InteractionCreate) {
	ms := discordgo.MessageSend{
		Embeds: i.Interaction.Message.Embeds,
	}

	// Only check for attachments if this evidence comes from a message,
	// the number of fields are different
	if i.Message.Embeds[0].Description == MessageEvidenceDescription {
		// Save attachments to the message because view links expire after 24h
		attachmentURLs := getAttachmentURLs(i.Interaction.Message.Embeds[0].Fields[6].Value)
		files := []*discordgo.File{}
		for _, attachmentURL := range attachmentURLs {
			client := http.Client{}
			resp, err := client.Get(attachmentURL)
			if err != nil {
				log.Warnf("error getting attachment from message (%s), error: %s", i.Message.ID, err)
				continue
			}
			defer resp.Body.Close()

			filename := path.Base(resp.Request.URL.Path)
			files = append(files, &discordgo.File{
				Name:        Spoiler + filename,
				ContentType: resp.Header.Get("content-type"),
				Reader:      resp.Body,
			})
		}

		// Add files to the attachment
		ms.Files = files
	}

	var cfg GuildConfig
	bot.DB.Where(&GuildConfig{ID: i.GuildID}).FirstOrCreate(&cfg)
	if cfg.EvidenceChannelSettingID == "" {
		err := bot.DG.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Title:   "Unable to submit report",
					Content: "Please use /" + Settings + " to set an evidence channel",
				},
			},
		)

		if err != nil {
			log.Errorf("error responding to settings interaction, err: %v", err)
		}
		return
	}

	ms.Embeds[0].Description = ""
	message, err := bot.DG.ChannelMessageSendComplex(cfg.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("unable to send message %w", err)
	} else {
		guild, err := bot.DG.Guild(i.GuildID)
		if err != nil {
			log.Warnf("unable to look up guild (%s), err: %s", i.GuildID, err)
			goto respond
		}
		// Save attachments to the message because view links expire after 24h
		userID := getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[0].Value)
		user, err := bot.DG.User(userID)
		if err != nil {
			log.Warnf("unable to look up user (%s, err: %v)", userID, err)
			goto respond
		}

		var notes string
		var previousReputation, currentReputation *int64
		// If notes exist, update. If not, add new
		for _, field := range i.Interaction.Message.Embeds[0].Fields {
			if field.Name == Notes {
				notes = Notes
			}
			if field.Name == CurrentReputation {
				value, _ := strconv.Atoi(field.Value)
				currentReputation = nullInt64(value)
			}
			if field.Name == PreviousReputation {
				value, _ := strconv.Atoi(field.Value)
				previousReputation = nullInt64(value)
			}
		}

		event := ModerationEvent{
			GuildId:            i.GuildID,
			GuildName:          guild.Name,
			UserID:             userID,
			UserName:           user.Username,
			ModeratedUserID:    i.GuildID + userID,
			Notes:              notes,
			PreviousReputation: previousReputation,
			CurrentReputation:  currentReputation,
			ModeratorID:        i.Interaction.Member.User.ID,
			ModeratorName:      i.Interaction.Member.User.Username,
			ReportURL: fmt.Sprintf(MessageURLTemplate,
				i.Interaction.GuildID,
				message.ChannelID,
				message.ID),
		}
		tx := bot.DB.Where(&ModeratedUser{ID: i.GuildID + userID}).Create(&event)
		tx.Commit()
	}

respond:
	err = bot.DG.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Submitted report: " +
					fmt.Sprintf(MessageURLTemplate,
						i.Interaction.GuildID,
						message.ChannelID,
						message.ID),
			},
		},
	)
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Gets user info from the `/query` command
func (bot *ModeratorBot) GetUserInfoFromChatCommandContext(i *discordgo.InteractionCreate) {
	userMentions := i.Interaction.ApplicationCommandData().Resolved.Users
	if len(userMentions) == 0 {
		log.Warn("user was not provided")
	}
	for userID := range userMentions {
		user := bot.GetModeratedUser(i.GuildID, userID)

		// TODO: Add more user information
		err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: GetUserInfoFromUserContext,
				Flags:    discordgo.MessageFlagsEphemeral,
				Content:  fmt.Sprintf("<@%s> has a reputation of %v", userID, *user.Reputation),
			},
		})
		if err != nil {
			log.Warn("error responding to interaction: %w", err)
		}
	}
}

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
	if len(message.Embeds) > 0 && len(message.Embeds[0].Fields) > 0 {
		fields = message.Embeds[0].Fields
		messageType = discordgo.InteractionResponseUpdateMessage
		authorID = getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[0].Value)
		for idx, field := range fields {
			if field.Name == CurrentReputation {
				user := bot.GetModeratedUser(i.GuildID, authorID)
				fields[idx].Value = fmt.Sprintf("%v", *user.Reputation)
			}
		}
	} else {
		messageContentNameField := OriginalMessageContent
		if len(message.Content) > 1024 {
			messageContentNameField = strings.Join([]string{OriginalMessageContent, "(truncated to 1024 characters)"}, " ")
		}
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
				Value:  fmt.Sprintf("%v", *user.Reputation),
				Inline: true,
			},
			{
				Name:   CurrentReputation,
				Value:  fmt.Sprintf("%v", *user.Reputation),
				Inline: true,
			},
			{
				Name:   "Link to original message",
				Value:  fmt.Sprintf(MessageURLTemplate, i.Interaction.GuildID, message.ChannelID, message.ID),
				Inline: true,
			},
			{
				Name:  messageContentNameField,
				Value: message.Content[:1024],
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
					Description: fmt.Sprintf(MessageEvidenceDescription, authorID),
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

// Returns a discordgo.InteractionResponse with an evidence report based on a message provided
func (bot *ModeratorBot) GenerateEvidenceReportFromUser(i *discordgo.InteractionCreate, user *discordgo.User) (resp *discordgo.InteractionResponse) {
	modUser := bot.GetModeratedUser(i.GuildID, user.ID)
	var fields []*discordgo.MessageEmbedField
	var messageType discordgo.InteractionResponseType
	var authorID string
	messageType = discordgo.InteractionResponseChannelMessageWithSource
	authorID = user.ID
	fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Author of message",
			Value:  fmt.Sprintf("<@%s>", user.ID),
			Inline: true,
		},
		{
			Name:   PreviousReputation,
			Value:  fmt.Sprintf("%v", *modUser.Reputation),
			Inline: true,
		},
		{
			Name:   CurrentReputation,
			Value:  fmt.Sprintf("%v", *modUser.Reputation),
			Inline: true,
		},
		{
			Name:   "Link to original message",
			Value:  "There is no original message as this evidence was collected by referencing the user directly.",
			Inline: true,
		},
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name: "Collected by",
		Value: fmt.Sprintf(`<@%s>
									%s (<t:%v:R>)`,
			i.Interaction.Member.User.ID,
			time.Now().UTC().Format(time.RFC1123Z),
			time.Now().Unix()),
	})

	return &discordgo.InteractionResponse{
		Type: messageType,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Evidence Report",
					Description: fmt.Sprintf(UserEvidenceDescription, authorID),
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

// Takes user-submitted notes and adds it to an in-progress report
func (bot *ModeratorBot) SaveEvidenceNotes(i *discordgo.InteractionCreate) {
	notes := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	// If notes exist, update. If not, add new
	for idx, field := range i.Interaction.Message.Embeds[0].Fields {
		if field.Name == Notes {
			i.Interaction.Message.Embeds[0].Fields[idx].Value = notes
			break
		}

		if idx == len(i.Interaction.Message.Embeds[0].Fields)-1 {
			i.Interaction.Message.Embeds[0].Fields = append(i.Interaction.Message.Embeds[0].Fields, &discordgo.MessageEmbedField{
				Name:  Notes,
				Value: notes,
			})
		}
	}

	err := bot.DG.InteractionRespond(i.Interaction,
		bot.GenerateEvidenceReportFromMessage(i, i.Interaction.Message))
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// User information and stats produced for the /query command and
// "Get info" when right clicking users
func (bot *ModeratorBot) userInfoIntegrationresponse(i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	user := i.Interaction.ApplicationCommandData().Resolved.Users[i.ApplicationCommandData().TargetID]

	if user.ID == "" {
		log.Warn("user was not provided")
		return &discordgo.InteractionResponseData{
			CustomID: GetUserInfoFromUserContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  "There was an error getting user info.",
		}
	}

	moderatedUser := bot.GetModeratedUser(i.GuildID, user.ID)
	return &discordgo.InteractionResponseData{
		CustomID: GetUserInfoFromUserContext,
		Flags:    discordgo.MessageFlagsEphemeral,
		Content:  fmt.Sprintf("<@%s> has a reputation of %v", user.ID, *moderatedUser.Reputation),
	}
}

// Returns a ModeratedUser record from the DB using server and user ID
// (a user can be in multiple servers)
func (bot *ModeratorBot) GetModeratedUser(GuildId string, userID string) (moderatedUser ModeratedUser) {
	guild, _ := bot.DG.Guild(GuildId)
	user, _ := bot.DG.User(userID)
	moderatedUser = ModeratedUser{
		UserName:   user.Username,
		UserID:     userID,
		GuildId:    GuildId,
		ID:         GuildId + userID,
		GuildName:  guild.Name,
		Reputation: nullInt64(1),
	}
	bot.DB.Where(&ModeratedUser{ID: GuildId + userID}).FirstOrCreate(&moderatedUser)
	return moderatedUser
}

// Updates a user reputation, given the source interaction and value to add
func (bot *ModeratorBot) ChangeUserReputation(i *discordgo.InteractionCreate, difference int) (err error) {
	userID := getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[0].Value)
	if userID == "" {
		return fmt.Errorf("unable to determine user ID from reference")
	}

	user, err := bot.DG.User(userID)
	if err != nil {
		return err
	}
	guild, err := bot.DG.Guild(i.GuildID)
	if err != nil {
		return err
	}

	userUpdate := ModeratedUser{
		UserID:    userID,
		UserName:  user.Username,
		GuildId:   i.GuildID,
		ID:        i.GuildID + userID,
		GuildName: guild.Name,
	}
	tx := bot.DB.Where(&ModeratedUser{ID: i.GuildID + userID}).FirstOrCreate(&userUpdate)
	if err != nil {
		return err
	}

	tx.Model(&ModeratedUser{}).Update("Reputation", *userUpdate.Reputation+int64(difference))
	return nil
}
