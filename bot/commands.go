package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	ChatCommands = []*discordgo.ApplicationCommand{
		{
			Name:        Help,
			Description: "How to use this bot",
		},
		{
			Name:        GetUserInfoFromChatCommandContext,
			Description: "Check a user's reputation information",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        UserOption,
					Description: "User to look up",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        CreatePoll,
			Description: "Create a poll",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        NumberOfPollChoices,
					Description: "Number of choices",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
					MaxValue:    10.0,
				},
			},
		},
		{
			Name:        Settings,
			Description: "Change settings",
		},
		{
			Name:        AddCustomSlashCommand,
			Description: "Create custom slash command",
		},
		{
			Name:        RemoveCustomSlashCommand,
			Description: "Remove custom slash command",
		},
	}
	UserCommands = []*discordgo.ApplicationCommand{
		{
			Name: GetUserInfoFromUserContext,
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: DocumentBehaviorFromUserContext,
			Type: discordgo.UserApplicationCommand,
		},
	}
	MessageCommands = []*discordgo.ApplicationCommand{
		{
			Name: GetUserInfoFromMessageContext,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: DocumentBehaviorFromMessageContext,
			Type: discordgo.MessageApplicationCommand,
		},
	}

	RegisteredCommands = []*discordgo.ApplicationCommand{}

	// This object is used to register chat commands, so if it's
	// not in here, it won't get registered properly. This is updated
	// during runtime.
	ConfiguredCommands = append(append(ChatCommands, UserCommands...), MessageCommands...)
)
