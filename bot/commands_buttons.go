package bot

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Updates a server setting according to the
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
		_ = bot.DG.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: bot.generalErrorDisplayedToTheUser(reason),
			},
		)
		return
	}

	_ = bot.DG.InteractionRespond(i.Interaction,
		bot.DocumentBehaviorFromMessage(i, &message))
}

func (bot *ModeratorBot) TakeEvidenceNotes(i *discordgo.InteractionCreate) {
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.SaveEvidenceNotes,
			Title:    "Add notes to this report",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    globals.EvidenceNotes,
						Label:       "Notes",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "Add notes such as reasoning or context",
						Required:    true,
						MinLength:   3,
						MaxLength:   500,
					},
				},
			}},
		},
	})
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
			log.Warn("error getting attachment from message (%s), error: %w", i.Message.ID, err)
			continue
		}
		defer resp.Body.Close()

		filename := path.Base(resp.Request.URL.Path)
		files = append(files, &discordgo.File{
			Name:        filename,
			ContentType: resp.Header.Get("content-type"),
			Reader:      resp.Body,
		})
	}

	// Add files to the attachment
	ms.Files = files

	// TODO: save event info
	sc := bot.getServerConfig(i.GuildID)
	message, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("Unable to send message %w", err)
	}

	_ = bot.DG.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Submitted report: " +
					fmt.Sprintf(globals.MessageURLTemplate,
						i.Interaction.GuildID,
						message.ChannelID,
						message.ID),
			},
		},
	)
}
