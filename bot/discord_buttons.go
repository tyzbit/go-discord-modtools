package bot

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Updates a server setting according to the
// column name (setting) and the value
func (bot *ModeratorBot) RespondToSettingsChoice(i *discordgo.InteractionCreate,
	setting string, value string) {

	tx := bot.DB.Model(&GuildConfig{}).
		Where(&GuildConfig{ID: i.Interaction.GuildID}).
		Update(setting, value)
	var interactionErr error

	if tx.RowsAffected != 1 {
		log.Debugf("unexpected number of rows affected updating guild settings: %v", tx.RowsAffected)
		interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: bot.generalErrorDisplayedToTheUser("Unable to save settings"),
		})
	} else {
		var cfg GuildConfig
		bot.DB.Where(&GuildConfig{ID: i.Interaction.GuildID}).First(&cfg)
		interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: bot.SettingsIntegrationResponse(cfg),
		})
	}

	if interactionErr != nil {
		log.Errorf("error responding to settings interaction, err: %v", interactionErr)
	}
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

	tx.Model(&ModeratedUser{}).Update("Reputation", *userUpdate.Reputation+int64(difference))
	return nil
}

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

// Deletes a custom command
func (bot *ModeratorBot) DeleteCustomSlashCommandFromButtonContext(i *discordgo.InteractionCreate, commandID string) {
	var cmds []CustomCommand
	bot.DB.Where(&CustomCommand{GuildConfigID: i.Interaction.GuildID}).Find(&cmds)

	var interactionErr error
	guild, _ := bot.DG.Guild(i.GuildID)
	if len(cmds) == 0 {
		log.Debugf("no registered commands returned for %s(%s)", guild.Name, guild.Name)
	}
	for _, cmd := range cmds {
		if cmd.ID == commandID {
			log.Infof("deleting custom command %s from server %s(%s)", commandID, guild.Name, guild.ID)
			tx := bot.DB.Where(&CustomCommand{GuildConfigID: guild.ID, ID: commandID}).Delete(&cmd)
			if tx.RowsAffected != 1 {
				log.Debugf("unexpected number of rows affected updating guild settings: %v", tx.RowsAffected)
				interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: bot.generalErrorDisplayedToTheUser("Unable to save settings"),
				})
			} else {
				defer bot.UpdateCommands()
				var cfg GuildConfig
				bot.DB.Where(&GuildConfig{ID: i.Interaction.GuildID}).First(&cfg)
				interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Title: "Command deleted",
						Flags: discordgo.MessageFlagsEphemeral,
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Command deleted",
								Description: fmt.Sprintf("/%s command deleted", cmd.Name),
								Color:       Green,
							},
						},
					},
				})
			}
		}
	}

	if interactionErr != nil {
		log.Errorf("error responding to settings interaction, err: %v", interactionErr)
	}
}

func (bot *ModeratorBot) PollUpdateHandler(i *discordgo.InteractionCreate) {
	mcd := i.Interaction.MessageComponentData()
	customID := mcd.CustomID
	guild, err := bot.DG.Guild(i.GuildID)
	if err != nil {
		log.Warnf("unable to look up guild (%s), err: %s", i.GuildID, err)
		return
	}

	var poll Poll
	var votes []Vote
	// Create if necessary
	bot.DB.Where(&Poll{GuildConfigID: guild.ID, ID: i.Message.ID}).Preload("Votes").FirstOrCreate(&poll)
	voteTx := bot.DB.Where(&Vote{GuildConfigID: guild.ID, PollID: i.Message.ID, CustomID: customID}).FirstOrCreate(&votes)

	for _, vote := range votes {
		if i.Member.User.ID == vote.UserID {
			err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: bot.generalInfoDisplayedToTheUser("You have already voted for that option"),
			})
			if err != nil {
				log.Warn("error responding to poll choice interaction: %w", err)
			}
			return
		}
	}

	log.Debugf("counting a vote for poll %s by %s in %s(%s)", poll.ID, i.Member.User.ID, guild.Name, guild.ID)
	poll.Votes = append(poll.Votes, Vote{
		PollID:        i.Message.ID,
		GuildConfigID: guild.ID,
		CustomID:      customID,
		UserID:        i.Member.User.ID,
	})
	voteTx.Updates(&votes)

	// TODO: update message with new tallies
	// TODO: show who voted for each option
}

// Starts a poll the user has customized
func (bot *ModeratorBot) EndPollFromButtonContext(i *discordgo.InteractionCreate) {
	// TODO: Update DB
	// TODO: Update message with results
	// TODO: Optionally notify user that created poll?
}
