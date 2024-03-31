package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	// Slash command
	Settings                     = "settings"
	SettingFailedResponseMessage = "Error changing setting"
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

// Sets setting choices from the `/settings` command
func (bot *ModeratorBot) SetSettingsFromChatCommandContext(i *discordgo.InteractionCreate) {
	log.Debug("handling settings request")
	if i.GuildID == "" {
		reason := "The bot does not have any per-user settings"
		// This is a DM, so settings cannot be changed
		err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bot.generalErrorDisplayedToTheUser(reason),
		})
		if err != nil {
			log.Errorf("error responding to settings DM"+Settings+", err: %v", err)
		}
		return
	} else {
		guild, err := bot.DG.Guild(i.Interaction.GuildID)
		if err != nil {
			guild.Name = "GuildLookupError"
		}

		var cfg GuildConfig
		bot.DB.Where(&GuildConfig{ID: i.GuildID}).First(&cfg)
		resp := bot.SettingsIntegrationResponse(cfg)
		err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: resp,
		})

		if err != nil {
			log.Errorf("error responding to slash command "+Settings+", err: %v", err)
		}
	}
}

// SettingsIntegrationResponse returns server settings in a *discordgo.InteractionResponseData
func (bot *ModeratorBot) SettingsIntegrationResponse(cfg GuildConfig) *discordgo.InteractionResponseData {
	channel, _ := bot.DG.Channel(cfg.EvidenceChannelSettingID)
	var evidenceChannelName, moderatorRoleName string
	if channel == nil {
		evidenceChannelName = "not set"
	} else {
		evidenceChannelName = "#" + channel.Name
	}
	moderatorRole, _ := bot.DG.State.Role(cfg.ID, cfg.ModeratorRoleSettingID)
	if moderatorRole == nil {
		moderatorRoleName = "not set"
	} else {
		moderatorRoleName = moderatorRole.Name
	}
	return &discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder:  EvidenceChannelSettingID + ": " + evidenceChannelName,
						MenuType:     discordgo.ChannelSelectMenu,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						CustomID:     EvidenceChannelSettingID,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						Placeholder: ModeratorRoleSettingID + ": " + moderatorRoleName,
						MenuType:    discordgo.RoleSelectMenu,
						CustomID:    ModeratorRoleSettingID,
					},
				},
			},
		},
	}
}
