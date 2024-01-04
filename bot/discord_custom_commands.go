package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// GetCustomCommandHandlers returns a map[string]func of command handlers for every ServerConfig
// TODO: filter out configured but not registered
func (bot *ModeratorBot) GetCustomCommandHandlers() (cmds map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	registeredGuildIds := []string{}
	cmds = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	bot.DB.Model(&GuildConfig{}).Pluck("guild_id", &registeredGuildIds)
	for _, regGuildId := range registeredGuildIds {
		var customCommands []CustomCommand
		bot.DB.Where(&CustomCommand{GuildConfigID: regGuildId}).Find(&customCommands)
		for _, customCommand := range customCommands {
			cmds[customCommand.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				bot.UseCustomSlashCommandFromChatCommandContext(i, customCommand.Content)
			}
		}
	}
	return cmds
}

// RegisterCustomCommandHandler registers all commands that are configured
// for a given ServerConfig
func (bot *ModeratorBot) RegisterCustomCommandHandler(cmds []CustomCommand) {
	commandsHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	commands, _ := bot.DG.ApplicationCommands("", "")
	for _, customCommand := range cmds {
		for _, registeredCommand := range commands {
			if customCommand.Name == registeredCommand.ID {
				log.Warnf("a saved server chat command conflicts with a global command and will be removed, %s",
					customCommand.Name)

				bot.DB.Model(&CustomCommand{}).Delete(CustomCommand{
					Name:          customCommand.Name,
					GuildConfigID: customCommand.GuildConfigID,
					Description:   customCommand.Description,
					Content:       customCommand.Content,
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

// UpdateCommands iterates through all configured commands and ensures
// they are registered, updated or removed
func (bot *ModeratorBot) UpdateCommands() (err error) {
	var commandsToCreate, commandsToEdit, commandsToDelete []*discordgo.ApplicationCommand

	// Get already-registered guild-specific commands from the database
	guildIds := []string{}
	bot.DB.Model(&GuildConfig{}).Pluck("guild_id", &guildIds)
	// Add an element for global commands (they do not have a Guild ID)
	guildIds = append(guildIds, "")
	for _, id := range guildIds {
		guildCommands, err := bot.DG.ApplicationCommands(bot.DG.State.User.ID, id)
		if err != nil {
			err = fmt.Errorf(", err: %w", err)
			log.Warnf("unable to look up server-specific commands for server %s", id)
			break
		}
		RegisteredCommands = append(RegisteredCommands, guildCommands...)

		var cfg GuildConfig
		bot.DB.Where(&GuildConfig{GuildID: id}).First(&cfg)

		// If these are guild commands, they won't be in ConfiguredCommands yet
		if id != "" {
			for _, configuredCommand := range cfg.CustomCommands {
				ConfiguredCommands = append(ConfiguredCommands, &discordgo.ApplicationCommand{
					Name:        configuredCommand.Name,
					Description: configuredCommand.Description,
					GuildID:     id,
				})
			}
		}
	}

	// Ensure configured commands get registered
	for _, configuredCommand := range ConfiguredCommands {
		// Compare against registered commands
		create := true
		for _, registeredCommand := range RegisteredCommands {
			nameMatch := configuredCommand.Name == registeredCommand.Name
			guildMatch := configuredCommand.GuildID == registeredCommand.GuildID
			descriptionMatch := configuredCommand.Description == registeredCommand.Description

			if nameMatch && guildMatch && !descriptionMatch {
				// This command is registered, but the description doesn't match
				registeredCommand.Description = configuredCommand.Description
				commandsToEdit = append(commandsToEdit, registeredCommand)
				break
			}

			if nameMatch && guildMatch && descriptionMatch {
				create = false
			}
		}

		if create {
			// If we're here, we've gone through all registered commands
			// but we didn't find this configured command, so
			// we need to create it
			commandsToCreate = append(commandsToCreate, configuredCommand)
		}
	}

	// Ensure extra registered commands get removed
	for _, registeredCommand := range RegisteredCommands {
		delete := true
		for _, configuredCommand := range ConfiguredCommands {
			nameMatch := configuredCommand.Name == registeredCommand.Name
			guildMatch := configuredCommand.GuildID == registeredCommand.GuildID

			if nameMatch && guildMatch {
				delete = false
				break
			}
		}

		if delete {
			commandsToDelete = append(commandsToDelete, registeredCommand)
		}

	}

	for _, command := range commandsToDelete {
		info := ""
		if command.GuildID != "" {
			guild, _ := bot.DG.Guild(command.GuildID)
			info = fmt.Sprintf(" from guild %s(%s)", guild.Name, command.GuildID)
		}
		log.Debugf("deleting command '/%s'", command.Name+info)
		err := bot.DG.ApplicationCommandDelete(bot.DG.State.User.ID,
			command.GuildID,
			command.ID)
		if err != nil {
			err = fmt.Errorf(", err: %w", err)
			log.Errorf("cannot remove command %s: %v", command.Name, err)
		}
	}

	for _, command := range commandsToEdit {
		info := ""
		if command.GuildID != "" {
			guild, _ := bot.DG.Guild(command.GuildID)
			info = fmt.Sprintf(" from guild %s(%s)", guild.Name, command.GuildID)
		}
		log.Debugf("editing command '/%s'", command.Name+info)
		editedCmd, err := bot.DG.ApplicationCommandEdit(bot.DG.State.User.ID,
			command.GuildID,
			command.ID,
			command)
		if err != nil {
			err = fmt.Errorf(", err: %w", err)
			log.Errorf("cannot update command '/%s': %v", command.Name, err)
		} else {
			RegisteredCommands = append(RegisteredCommands, editedCmd)
		}
	}

	for _, command := range commandsToCreate {
		info := ""
		if command.GuildID != "" {
			guild, _ := bot.DG.Guild(command.GuildID)
			info = fmt.Sprintf(" from guild %s(%s)", guild.Name, command.GuildID)
		}
		log.Debugf("creating command '/%s'", command.Name+info)
		newCmd, err := bot.DG.ApplicationCommandCreate(bot.DG.State.User.ID,
			command.GuildID,
			command)
		if err != nil {
			err = fmt.Errorf(", err: %w", err)
			log.Errorf("cannot create command '/%s': %v", command.Name, err)
		} else {
			RegisteredCommands = append(RegisteredCommands, newCmd)
		}
	}
	return err
}
