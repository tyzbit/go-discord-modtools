package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// This produces the help text seen from the chat commant `/help`
func (bot *ModeratorBot) GetHelpFromChatCommandContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Moder8or Bot Help",
					Description: globals.BotHelpText,
					Color:       globals.Purple,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("error responding to help command "+globals.Help+", err: %v", err)
	}
}

// Produces bot stats, server-specific if called in a server and
// a total summary if DMed by a configured administrator
func (bot *ModeratorBot) GetStatsFromChatCommandContext(i *discordgo.InteractionCreate) {
	directMessage := (i.GuildID == "")
	var stats botStats
	logMessage := ""
	if !directMessage {
		log.Debug("handling stats request")
		stats = bot.getServerStats(i.GuildID)
		guild, err := bot.DG.Guild(i.GuildID)
		if err != nil {
			log.Errorf("unable to look up server by id: %v", i.GuildID+", "+fmt.Sprintf("%v", err))
			return
		}
		logMessage = "sending stats response to " + i.Member.User.Username + "(" + i.Member.User.ID + ") in " +
			guild.Name + "(" + guild.ID + ")"
	} else {
		log.Debug("handling stats DM request")
		// We can be sure now the request was a direct message
		// Deny by default
		administrator := false

	out:
		// TODO: allow adding, removing and looking up admins in the DB
		for _, id := range bot.Config.AdminIds {
			if i.User.ID == id {
				administrator = true

				// This prevents us from checking all IDs now that
				// we found a match but is a fairly ineffectual
				// optimization since config.AdminIds will probably
				// only have dozens of IDs at most
				break out
			}
		}

		if !administrator {
			log.Errorf("did not respond to global stats command from %v(%v), because user is not an administrator",
				i.User.Username, i.User.ID)

			err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Stats are not available in DMs",
							Color: globals.Purple,
						},
					},
				},
			})
			if err != nil {
				log.Errorf("error responding to slash command "+globals.Stats+", err: %v", err)
			}
			return
		}
		stats = bot.getGlobalStats()
		logMessage = "sending global " + globals.Stats + " response to " + i.User.Username + "(" + i.User.ID + ")"
	}

	log.Info(logMessage)

	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  "🏛️ Moder8or Bot Stats",
					Fields: structToPrettyDiscordFields(stats, directMessage),
					Color:  globals.Purple,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("error responding to slash command "+globals.Stats+", err: %v", err)
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
			log.Errorf("error responding to settings DM"+globals.Settings+", err: %v", err)
		}
		return
	} else {
		guild, err := bot.DG.Guild(i.Interaction.GuildID)
		if err != nil {
			guild.Name = "GuildLookupError"
		}

		sc := bot.getServerConfig(i.GuildID)
		resp := bot.SettingsIntegrationResponse(sc)
		err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: resp,
		})

		if err != nil {
			log.Errorf("error responding to slash command "+globals.Settings+", err: %v", err)
		}
	}
}

// Gets user info from the `/query` command
func (bot *ModeratorBot) GetUserInfoFromChatCommandContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User.ID == "" {
		log.Warn("user was not provided")
	}

	user := ModeratedUser{}
	bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: i.Interaction.Member.User.ID}).First(&user)

	// TODO: Add more user information
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.GetUserInfoFromUserContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, user.Reputation.Int64),
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Creates a new slash command
func (bot *ModeratorBot) ConfigureCustomSlashCommandFromChatCommandContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.SaveCustomSlashCommand,
			Title:    "Create new simple custom slash command",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  globals.CustomSlashName,
						Label:     globals.CustomSlashName,
						Style:     discordgo.TextInputShort,
						Required:  true,
						MinLength: 1,
						MaxLength: 32,
					},
				},
			},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  globals.CustomSlashCommandDescription,
							Label:     globals.CustomSlashCommandDescription,
							Style:     discordgo.TextInputShort,
							Required:  true,
							MinLength: 1,
							MaxLength: 32,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.CustomSlashCommandContent,
							Label:       globals.CustomSlashName,
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Text this slash command sends to the channel in which it is called",
							Required:    true,
							MinLength:   1,
							MaxLength:   globals.MaxMessageContentLength,
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Errorf("error showing custom slash command creation modal, err: %v", err)
	}
}

// Creates a new slash command
func (bot *ModeratorBot) UseCustomSlashCommandFromChatCommandContext(i *discordgo.InteractionCreate, content string) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content}})
	if err != nil {
		log.Errorf("error responding to use of custom command, err: %v", err)
	}
}
