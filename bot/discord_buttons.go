package bot

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Updates a server setting according to the
// column name (setting) and the value
func (bot *ModeratorBot) RespondToSettingsChoice(i *discordgo.InteractionCreate,
	setting string, value interface{}) {
	sc, ok := bot.updateServerSetting(i.Interaction.GuildID, setting, value)
	var interactionErr error

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

// Updates a user reputation, given the source interaction and
// whether to increase (TRUE) or decrease (FALSE)
func (bot *ModeratorBot) ChangeUserReputation(i *discordgo.InteractionCreate, increase bool) {
	userID := getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[1].Value)
	user := ModeratedUser{}
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: userID}).First(&user)

	if tx.RowsAffected > 1 {
		log.Errorf("unexpected number of rows affected getting user reputation: %v", tx.RowsAffected)
		return
	} else if tx.RowsAffected == 0 {
		user.UserID = userID
	}

	if increase {
		user.Reputation = sql.NullInt64{
			Valid: true,
			Int64: user.Reputation.Int64 + 1,
		}
	} else {
		user.Reputation = sql.NullInt64{
			Valid: true,
			Int64: user.Reputation.Int64 - 1,
		}
	}

	err := bot.UpdateModeratedUser(user)
	if err != nil {
		log.Warn("unable to update user moderation record, err: %w", err)
	}
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

	// TODO: save event info
	sc := bot.getServerConfig(i.GuildID)
	if sc.EvidenceChannelSettingID == "" {
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
	message, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("unable to send message %w", err)
	} else {
		guild, err := bot.DG.Guild(i.GuildID)
		if err != nil {
			log.Warnf("unable to look up guild (%s), err: %s", i.GuildID, err)
			goto respond
		}
		// Save attachments to the message because view links expire after 24h
		userID := getUserIDFromDiscordReference(i.Interaction.Message.Embeds[0].Fields[1].Value)
		user, err := bot.DG.User(userID)
		if err != nil {
			log.Warnf("unable to look up user (%s, err: %v)", userID, err)
			goto respond
		}

		var notes string
		var previousReputation, currentReputation sql.NullInt64
		// If notes exist, update. If not, add new
		for _, field := range i.Interaction.Message.Embeds[0].Fields {
			if field.Name == Notes {
				notes = Notes
			}
			if field.Name == CurrentReputation {
				value, _ := strconv.Atoi(field.Value)
				currentReputation = sql.NullInt64{
					Valid: true,
					Int64: int64(value),
				}
			}
			if field.Name == PreviousReputation {
				value, _ := strconv.Atoi(field.Value)
				previousReputation = sql.NullInt64{
					Valid: true,
					Int64: int64(value),
				}
			}
		}

		bot.createModerationEvent(ModerationEvent{
			UUID:               uuid.New().String(),
			ServerID:           i.GuildID,
			ServerName:         guild.Name,
			UserID:             userID,
			UserName:           user.Username,
			Notes:              notes,
			PreviousReputation: previousReputation,
			CurrentReputation:  currentReputation,
			ModeratorID:        i.Interaction.Member.User.ID,
			ModeratorName:      i.Interaction.Member.User.Username,
			ReportURL: fmt.Sprintf(MessageURLTemplate,
				i.Interaction.GuildID,
				message.ChannelID,
				message.ID),
		})
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
