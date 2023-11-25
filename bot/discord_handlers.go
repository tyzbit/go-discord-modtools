package bot

import (
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
	// bot.DeleteAllCommands()
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
	// Technically app actions are commands too, but those are in commands_message.go and commands_user.go
	// We don't pass the session to these because you can get that from bot.DG
	commandsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Chat commands (slash commands)
		globals.Help: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetHelpFromChatCommandContext(i)
		},
		// Stats does not create an InteractionEvent
		globals.Stats: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetStatsFromChatCommandContext(i)
		},
		globals.Settings: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.SetSettingsFromChatCommandContext(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.GetUserInfoFromChatCommandContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromChatCommandContext(i)
		},
		// Message actions (right click or long-press message)
		globals.DocumentBehaviorFromUserContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DocumentBehaviorFromUserContext(i)
		},

		// User actions (right click or long-press user)
		// TODO: error will be handled once the functions are ready
		globals.GetUserInfoFromUserContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromUserContext(i)
		},

		// Message actions (right click or long-press message)
		globals.DocumentBehaviorFromMessageContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DocumentBehaviorFromMessageContext(i)
		},
		// TODO: error will be handled once the functions are ready
		globals.GetUserInfoFromMessageContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromMessageContext(i)
		},
	}

	buttonHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Settings buttons/choices
		globals.NotifyWhenReputationIsBelowSetting: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.RespondToSettingsChoice(i, "notify_when_reputation_is_below_setting", mcd.Values[0])
		},
		globals.NotifyWhenReputationIsAboveSetting: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.RespondToSettingsChoice(i, "notify_when_reputation_is_above_setting", mcd.Values[0])
		},
		globals.EvidenceChannelSettingID: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.RespondToSettingsChoice(i, "evidence_channel_setting_id", mcd.Values[0])
		},
		globals.IncreaseUserReputation: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ChangeUserReputation(i, true)
			// TODO: edit the original message we posted instead of posting a new one
			bot.GetUserInfoFromMessageContext(i)
		},
		globals.DecreaseUserReputation: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ChangeUserReputation(i, false)
			// TODO: edit the original message we posted instead of posting a new one
			bot.GetUserInfoFromMessageContext(i)
		},
	}

	// modalHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	// 	globals.DocumentBehaviorFromMessageContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 		bot.DocumentBehaviorFromModalSubmission(i)
	// 	},
	// }

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := buttonHandlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
		// case discordgo.InteractionModalSubmit:
		// 	if h, ok := modalHandlers[i.ModalSubmitData().CustomID]; ok {
		// 		h(s, i)
		// 	}
	}
}
