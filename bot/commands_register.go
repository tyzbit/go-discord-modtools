package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

func (bot *ModeratorBot) UpdateCommands() {
	var commandsToCreate, commandsToEdit, commandsToDelete []*discordgo.ApplicationCommand

	// Get already-registered guild-specific commands from the database
	guildIds := []string{}
	bot.DB.Model(&ServerRegistration{}).Pluck("discord_id", &guildIds)
	// Add an element for global commands (they do not have a Guild ID)
	guildIds = append(guildIds, "")
	for _, id := range guildIds {
		guildCommands, err := bot.DG.ApplicationCommands(bot.DG.State.User.ID, id)
		if err != nil {
			log.Warnf("unable to look up server-specific commands for server %s", id)
			break
		}
		globals.RegisteredCommands = append(globals.RegisteredCommands, guildCommands...)

		sc := bot.getServerConfig(id)

		// If these are guild commands, they won't be in globals.ConfiguredCommands yet
		if id != "" {
			for _, configuredCommand := range sc.CustomCommands {
				globals.ConfiguredCommands = append(globals.ConfiguredCommands, &discordgo.ApplicationCommand{
					Name:        configuredCommand.Name,
					Description: configuredCommand.Description,
					GuildID:     id,
				})
			}
		}
	}

	// Ensure configured commands get registered
	for _, configuredCommand := range globals.ConfiguredCommands {
		// Compare against registered commands
		create := true
		for _, registeredCommand := range globals.RegisteredCommands {
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
	for _, registeredCommand := range globals.RegisteredCommands {
		delete := true
		for _, configuredCommand := range globals.ConfiguredCommands {
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
			log.Errorf("cannot update command '/%s': %v", command.Name, err)
		} else {
			globals.RegisteredCommands = append(globals.RegisteredCommands, editedCmd)
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
			log.Errorf("cannot create command %s: %v", command.Name, err)
		} else {
			globals.RegisteredCommands = append(globals.RegisteredCommands, newCmd)
		}
	}
}
