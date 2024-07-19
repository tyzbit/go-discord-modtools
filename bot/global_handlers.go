package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// BotReadyHandler is called when the bot is considered ready to use the Discord session
func (bot *ModeratorBot) BotReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	// r.Guilds has all of our connected servers, so we should
	// update server registrations and set any registered servers
	// not in r.Guilds as inactive
	bot.UpdateInactiveRegistrations(r.Guilds)

	err := bot.updateServersWatched()
	if err != nil {
		log.Error("unable to update servers watched")
	}
}

// GuildCreateHandler is called whenever the bot joins a new guild.
func (bot *ModeratorBot) GuildCreateHandler(s *discordgo.Session, gc *discordgo.GuildCreate) {
	if gc.Guild.Unavailable {
		return
	}

	err := bot.registerOrUpdateServer(gc.Guild)
	if err != nil {
		log.Errorf("unable to register or update server: %v", err)
	}

	// Start watching RSS feeds
	go bot.StartWatchingRegisteredFeeds(gc.Guild.ID)
}

// GuildDeleteHandler is called whenever the bot leaves a guild.
func (bot *ModeratorBot) GuildDeleteHandler(s *discordgo.Session, gd *discordgo.GuildDelete) {
	if gd.Guild.Unavailable {
		return
	}

	log.Infof("guild %s(%s) deleted (bot was probably kicked)", gd.Guild.Name, gd.Guild.ID)
	err := bot.registerOrUpdateServer(gd.BeforeDelete)
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
		Help: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetHelpFromChatCommandContext(i)
		},
		Settings: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.SetSettingsFromChatCommandContext(i)
		},
		GetUserInfoFromChatCommandContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromChatCommandContext(i)
		},
		// Message actions (right click or long-press message)
		DocumentBehaviorFromUserContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DocumentBehaviorFromUserContext(i)
		},

		// User actions (right click or long-press user)
		GetUserInfoFromUserContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromUserContext(i)
		},

		// Message actions (right click or long-press message)
		DocumentBehaviorFromMessageContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DocumentBehaviorFromMessageContext(i)
		},
		GetUserInfoFromMessageContext: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.GetUserInfoFromMessageContext(i)
		},
		AddCustomSlashCommand: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.CreateCustomSlashCommandFromChatCommandContext(i)
		},
		RemoveCustomSlashCommand: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DeleteCustomSlashCommandFromChatCommandContext(i)
		},
		ConfigureRSSFeed: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ConfigureRSSFeedFromChatCommandContext(i)
		},
		ListRSSFeeds: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ListRSSFeedsFromChatCommandContext(i)
		},
		SelectRSSFeedForDeletion: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.DeleteRSSFeedFromChatCommandContext(i)
		},
	}

	buttonHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// Settings buttons/choices
		EvidenceChannelSettingID: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.RespondToSettingsChoice(i, "evidence_channel_setting_id", string(mcd.Values[0]))
		},
		ModeratorRoleSettingID: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			botGuildMember, err := bot.DG.GuildMember(i.GuildID, i.Interaction.Message.Author.ID)
			if err != nil {
				log.Warn("Unable to look up bot in list of guild members, err: %w", err)

			} else {
				for idx, role := range botGuildMember.Roles {
					if mcd.Values[0] == role {
						err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: bot.generalErrorDisplayedToTheUser("Select a different role, this bot role cannot be used"),
						})
						if err != nil {
							log.Warn("error responding to interaction: %w", err)
						}
						break
					}
					if idx == len(botGuildMember.Roles)-1 {
						bot.RespondToSettingsChoice(i, "moderator_role_setting_id", string(mcd.Values[0]))
					}
				}
			}
		},
		IncreaseUserReputation: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ChangeUserReputation(i, 1)
			bot.DocumentBehaviorFromButtonContext(i)
		},
		DecreaseUserReputation: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ChangeUserReputation(i, -1)
			bot.DocumentBehaviorFromButtonContext(i)
		},
		ShowEvidenceCollectionModal: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ShowEvidenceCollectionModal(i)
		},
		SubmitReport: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.SubmitReport(i)
		},
		DeleteCustomSlashCommand: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mcd := i.MessageComponentData()
			bot.DeleteCustomSlashCommandFromButtonContext(i, mcd.Values[0])
		},
	}

	modalHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		SaveEvidenceNotes: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.SaveEvidenceNotes(i)
		},
		SaveCustomSlashCommand: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.SaveCustomSlashCommand(i)
		},
		ConfigureRSSFeed: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.ConfigureRSSFeed(i)
		},
	}

	customCommandHandlers := bot.GetCustomCommandHandlers()
	var cfg GuildConfig
	bot.DB.Where(&GuildConfig{ID: i.GuildID}).First(&cfg)
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			if bot.isAllowed(cfg, i.Member) {
				h(s, i)
			} else {
				err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: bot.permissionsErrorDisplayedToTheUser(),
				})
				if err != nil {
					log.Warn("error responding to application command interaction: %w", err)
				}
			}
		}
		if h, ok := customCommandHandlers[i.ApplicationCommandData().Name]; ok {
			if bot.isAllowed(cfg, i.Member) {
				h(s, i)
			} else {
				err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: bot.permissionsErrorDisplayedToTheUser(),
				})
				if err != nil {
					log.Warn("error responding to custom command interaction: %w", err)
				}
			}
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := buttonHandlers[i.MessageComponentData().CustomID]; ok {
			if bot.isAllowed(cfg, i.Member) {
				h(s, i)
			} else {
				err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: bot.permissionsErrorDisplayedToTheUser(),
				})
				if err != nil {
					log.Warn("error responding to interaction message interaction: %w", err)
				}
			}
		}
	case discordgo.InteractionModalSubmit:
		if h, ok := modalHandlers[i.ModalSubmitData().CustomID]; ok {
			if bot.isAllowed(cfg, i.Member) {
				h(s, i)
			} else {
				err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: bot.permissionsErrorDisplayedToTheUser(),
				})
				if err != nil {
					log.Warn("error responding to modal submit interaction: %w", err)
				}
			}
		}
	}
}
