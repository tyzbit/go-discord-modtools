package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	// Chat commands
	GetUserInfoFromChatCommandContext = "query"
	AddCustomSlashCommand             = "addcommand"
	RemoveCustomSlashCommand          = "deletecommand"

	// Modals
	SaveCustomSlashCommand   = "Save custom slash command"
	DeleteCustomSlashCommand = "Remove custom slash command"

	// Modal options
	CustomSlashName               = "Name for this custom slash command"
	CustomSlashCommandDescription = "Description for this custom slash command"
	CustomSlashCommandContent     = "Message to paste if this command is used"
	CustomCommandIdentifier       = "Custom command: "

	// Constants
	MaxCommandContentLength     = 32   // https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
	MaxMessageContentLength     = 2000 // https://discord.com/developers/docs/resources/channel#create-message
	MaxDescriptionContentLength = 100  // https://discord.com/developers/docs/interactions/application-commands#application-command-object
)

// GetCustomCommandHandlers returns a map[string]func of command handlers for every ServerConfig
func (bot *ModeratorBot) GetCustomCommandHandlers() (cmds map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	activeGuildIds := []string{}
	cmds = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	bot.DB.Model(&GuildConfig{}).Where(&GuildConfig{Active: nullBool(true)}).Pluck("id", &activeGuildIds)
	for _, regGuildId := range activeGuildIds {
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

// RegisterCustomCommandHandler registers the provided commands
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
	bot.DB.Model(&GuildConfig{}).Pluck("id", &guildIds)
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
		bot.DB.Where(&GuildConfig{ID: id}).First(&cfg)

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

// Deletes a custom command
func (bot *ModeratorBot) DeleteCustomSlashCommandFromButtonContext(i *discordgo.InteractionCreate, commandID string) {
	var cmds []CustomCommand
	bot.DB.Where(&CustomCommand{GuildConfigID: i.Interaction.GuildID}).Find(&cmds)

	var interactionErr error
	guild, _ := bot.DG.Guild(i.GuildID)
	if len(cmds) == 0 {
		log.Debugf("no registered commands returned for %s(%s)", guild.Name, guild.Name)
	}
	for _, cmd := range cmds {
		if cmd.ID == commandID {
			log.Infof("deleting custom command %s from server %s(%s)", commandID, guild.Name, guild.ID)
			tx := bot.DB.Where(&CustomCommand{GuildConfigID: guild.ID, ID: commandID}).Delete(&cmd)
			if tx.RowsAffected != 1 {
				log.Debugf("unexpected number of rows affected updating guild settings: %v", tx.RowsAffected)
				interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: bot.generalErrorDisplayedToTheUser("Unable to save settings"),
				})
			} else {
				defer bot.UpdateCommands()
				var cfg GuildConfig
				bot.DB.Where(&GuildConfig{ID: i.Interaction.GuildID}).First(&cfg)
				interactionErr = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Title: "Command deleted",
						Flags: discordgo.MessageFlagsEphemeral,
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Command deleted",
								Description: fmt.Sprintf("/%s command deleted", cmd.Name),
								Color:       Green,
							},
						},
					},
				})
			}
		}
	}

	if interactionErr != nil {
		log.Errorf("error responding to settings interaction, err: %v", interactionErr)
	}
}

// Creates a new slash command
func (bot *ModeratorBot) CreateCustomSlashCommandFromChatCommandContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: SaveCustomSlashCommand,
			Title:    "Create new simple custom slash command",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  CustomSlashName,
						Label:     CustomSlashName,
						Style:     discordgo.TextInputShort,
						Required:  true,
						MinLength: 1,
						MaxLength: MaxCommandContentLength,
					},
				},
			},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  CustomSlashCommandDescription,
							Label:     CustomSlashCommandDescription,
							Style:     discordgo.TextInputShort,
							Required:  true,
							MinLength: 1,
							MaxLength: MaxDescriptionContentLength - len(CustomCommandIdentifier),
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    CustomSlashCommandContent,
							Label:       CustomSlashName,
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Text this slash command sends to the channel in which it is called",
							Required:    true,
							MinLength:   1,
							MaxLength:   MaxMessageContentLength,
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

func (bot *ModeratorBot) DeleteCustomSlashCommandFromChatCommandContext(i *discordgo.InteractionCreate) {
	var cmds []CustomCommand
	bot.DB.Where(&CustomCommand{GuildConfigID: i.GuildID}).Find(&cmds)

	var options []discordgo.SelectMenuOption
	for _, cmd := range cmds {
		options = append(options, discordgo.SelectMenuOption{
			Label:       fmt.Sprintf("/%s", cmd.Name),
			Description: strings.ReplaceAll(cmd.Description, CustomCommandIdentifier, ""),
			Value:       cmd.ID,
		})
	}

	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder: "Command to delete",
							CustomID:    DeleteCustomSlashCommand,
							Options:     options,
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

// Takes user-submitted command and adds it to the server config and registers
// it in the guild
func (bot *ModeratorBot) SaveCustomSlashCommand(i *discordgo.InteractionCreate) {
	guild, _ := bot.DG.Guild(i.GuildID)
	name := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	description := i.Interaction.ModalSubmitData().
		Components[1].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	content := i.Interaction.ModalSubmitData().
		Components[2].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	customCommand := CustomCommand{
		GuildConfigID: guild.ID,
		Name:          strings.ToLower(name),
		Description:   CustomCommandIdentifier + description,
		Content:       content,
	}

	ird := discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("Custom command /%s created", customCommand.Name),
				Description: fmt.Sprintf(`Description:
				%s
				
				Content:
					%s`, description, customCommand.Content),
				Color: Purple,
			},
		},
	}

	if len(customCommand.Description) > MaxDescriptionContentLength {
		ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Please limit the description to %v", MaxDescriptionContentLength-len(CustomCommandIdentifier)))
	} else if len(customCommand.Content) > MaxMessageContentLength {
		ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Please limit the message length to %v", MaxMessageContentLength))
	} else if strings.Contains(customCommand.Name, " ") {
		ird = *bot.generalErrorDisplayedToTheUser("Command names may not have spaces")
	} else {
		info := fmt.Sprintf(" from guild %s(%s)", guild.Name, guild.ID)
		log.Debugf("creating command '/%s'", name+info)
		cmd, err := bot.DG.ApplicationCommandCreate(bot.DG.State.User.ID,
			guild.ID,
			&discordgo.ApplicationCommand{
				Name:        customCommand.Name,
				Description: customCommand.Description,
				GuildID:     guild.ID,
			})
		if err != nil {
			ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Unable to create command, err: %v", err))
		}

		customCommand.ID = cmd.ID
		bot.DB.Model(&CustomCommand{}).Where(&CustomCommand{GuildConfigID: guild.ID}).Create(&customCommand)

		// Register commands with Discord API
		bot.RegisterCustomCommandHandler([]CustomCommand{customCommand})
	}

	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &ird,
	})
	if err != nil {
		log.Errorf("error responding to custom slash command creation, err: %v", err)
	}
}
