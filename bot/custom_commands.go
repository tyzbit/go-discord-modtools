package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (bot *ModeratorBot) GetCustomCommandHandlers() (cmds map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	registeredServerIDs := []string{}
	cmds = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	bot.DB.Model(&ServerRegistration{}).Pluck("discord_id", &registeredServerIDs)
	for _, regServerId := range registeredServerIDs {
		sc := bot.getServerConfig(regServerId)
		for _, customCommand := range sc.CustomCommands {
			cmds[customCommand.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				bot.UseCustomSlashCommandFromChatCommandContext(i, customCommand.Content)
			}
		}
	}
	return cmds
}

func (bot *ModeratorBot) RegisterCustomCommandHandler(sc ServerConfig) {
	commandsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	commands, _ := bot.DG.ApplicationCommands("", "")
	for _, customCommand := range sc.CustomCommands {
		for _, registeredCommand := range commands {
			if customCommand.Name == registeredCommand.ID {
				log.Warnf("a saved server chat command conflicts with a global command and will be removed, %s",
					customCommand.Name)

				bot.DB.Model(&CustomCommand{}).Delete(CustomCommand{
					Name:        customCommand.Name,
					DiscordId:   customCommand.DiscordId,
					Description: customCommand.Description,
					Content:     customCommand.Content,
				})
				// I don't think this is how you do this but I
				// can't remember the right way right now lol
				goto skip_custom_command
			}
		}

		commandsHandlers[customCommand.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.UseCustomSlashCommandFromChatCommandContext(i, customCommand.Content)
		}
	skip_custom_command:
	}
}
