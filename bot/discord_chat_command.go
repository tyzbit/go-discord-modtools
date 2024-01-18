package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
					Description: BotHelpText,
					Color:       Purple,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("error responding to help command "+Help+", err: %v", err)
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

// Gets user info from the `/query` command
func (bot *ModeratorBot) GetUserInfoFromChatCommandContext(i *discordgo.InteractionCreate) {
	if i.Interaction.Member.User.ID == "" {
		log.Warn("user was not provided")
	}

	user := ModeratedUser{}
	bot.DB.Where(&ModeratedUser{ID: i.GuildID + i.Interaction.Member.User.ID}).First(&user)

	// TODO: Add more user information
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: GetUserInfoFromUserContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  fmt.Sprintf("<@%s> has a reputation of %v", i.Interaction.Member.User.ID, &user.Reputation),
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
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
