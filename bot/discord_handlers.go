package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// BotReadyHandler is called when the bot is considered ready to use the Discord session
func (bot *ModeratorBot) BotReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	// r.Guilds has all of our connected servers, so we should
	// update server registrations and set any registered servers
	// not in r.Guilds as inactive
	bot.updateInactiveRegistrations(r.Guilds)

	// Use this to clean up commands if IDs have changed
	// TODO remove later if unnecessary
	// log.Debug("removing all commands")
	// bot.deleteAllCommands()
	// var err error
	// globals.RegisteredCommands, err = bot.DG.ApplicationCommandBulkOverwrite(bot.DG.State.User.ID, "", globals.Commands)
	log.Debug("registering slash commands")
	registeredCommands, err := bot.DG.ApplicationCommands(bot.DG.State.User.ID, "")
	if err != nil {
		log.Errorf("unable to look up registered application commands, err: %s", err)
	} else {
		for _, botCommand := range globals.Commands {
			for i, registeredCommand := range registeredCommands {
				// Check if this registered command matches a configured bot command
				if botCommand.Name == registeredCommand.Name {
					// Only update if it differs from what's already registered
					if botCommand != registeredCommand {
						editedCmd, err := bot.DG.ApplicationCommandEdit(bot.DG.State.User.ID, "", registeredCommand.ID, botCommand)
						if err != nil {
							log.Errorf("cannot update command %s: %v", botCommand.Name, err)
						}
						globals.RegisteredCommands = append(globals.RegisteredCommands, editedCmd)

						// Bot command was updated, so skip to the next bot command
						break
					}
				}

				// Check on the last item of registeredCommands
				if i == len(registeredCommands) {
					// This is a stale registeredCommand, so we should delete it
					err := bot.DG.ApplicationCommandDelete(bot.DG.State.User.ID, "", registeredCommand.ID)
					if err != nil {
						log.Errorf("cannot remove command %s: %v", registeredCommand.Name, err)
					}
				}
			}

			// If we're here, then we have a command that needs to be registered
			createdCmd, err := bot.DG.ApplicationCommandCreate(bot.DG.State.User.ID, "", botCommand)
			if err != nil {
				log.Errorf("cannot update command %s: %v", botCommand.Name, err)
			}
			globals.RegisteredCommands = append(globals.RegisteredCommands, createdCmd)
			if err != nil {
				log.Errorf("cannot update commands: %v", err)
			}
		}
	}

	err = bot.updateServersWatched()
	if err != nil {
		log.Error("unable to update servers watched")
	}
}

// GuildCreateHandler is called whenever the bot joins a new guild.
func (bot *ModeratorBot) GuildCreateHandler(s *discordgo.Session, gc *discordgo.GuildCreate) {
	if gc.Guild.Unavailable {
		return
	}

	err := bot.registerOrUpdateServer(gc.Guild, false)
	if err != nil {
		log.Errorf("unable to register or update server: %v", err)
	}
}

// GuildDeleteHandler is called whenever the bot leaves a guild.
func (bot *ModeratorBot) GuildDeleteHandler(s *discordgo.Session, gd *discordgo.GuildDelete) {
	if gd.Guild.Unavailable {
		return
	}

	log.Infof("guild %s(%s) deleted (bot was probably kicked)", gd.Guild.Name, gd.Guild.ID)
	err := bot.registerOrUpdateServer(gd.BeforeDelete, true)
	if err != nil {
		log.Errorf("unable to register or update server: %v", err)
	}
}

// InteractionInit configures all interactive commands
func (bot *ModeratorBot) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		globals.Help: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				log.Errorf("error responding to hexlp command "+globals.Help+", err: %v", err)
			}
		},
		// bot.createModerationEvent can handle both the moderate slash command and the app menu function
		// TODO: error will be handled once the functions are ready
		globals.ModerateUsingUser: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ModerateMenuFromUser(i)
		},
		globals.ModerateUsingMessage: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ModerateMenuFromMessage(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.Query: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.CheckUserReputation(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.CheckUserReputation: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.CheckUserReputation(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.CheckUserReputationUsingMessage: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.CheckUserReputation(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.CollectMessageAsEvidence: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.CollectMessageAsEvidence(i)
		},
		// Stats does not create an InteractionEvent
		globals.Stats: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		globals.Settings: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Debug("handling settings request")
			if i.GuildID == "" {
				// This is a DM, so settings cannot be changed
				err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: bot.settingsDMFailureIntegrationResponse(),
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

				bot.createInteractionEvent(InteractionEvent{
					UserID:        i.Interaction.Member.User.ID,
					Username:      i.Interaction.Member.User.Username,
					InteractionId: i.ID,
					ChannelId:     i.Interaction.ChannelID,
					ServerID:      i.Interaction.GuildID,
					ServerName:    guild.Name,
				})

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
		},
	}

	buttonHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Settings buttons/choices
		globals.NotifyWhenReputationIsBelowSetting: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.respondToSettingsChoice(i, "notify_when_reputation_is_below_setting", mcd.Values[0])
		},
		globals.NotifyWhenReputationIsAboveSetting: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.respondToSettingsChoice(i, "notify_when_reputation_is_above_setting", mcd.Values[0])
		},
		globals.EvidenceChannelSettingID: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.respondToSettingsChoice(i, "evidence_channel_setting_id", mcd.Values[0])
		},
	}

	modalHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		globals.ModerateUsingUser: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ModerateActionFromUser(i)
		},
		globals.ModerateUsingMessage: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ModerateActionFromMessage(i)
		},
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := buttonHandlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
	case discordgo.InteractionModalSubmit:
		if h, ok := modalHandlers[i.ModalSubmitData().CustomID]; ok {
			h(s, i)
		}
	}
}
