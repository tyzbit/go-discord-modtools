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
		{
			Name:        ConfigureRSSFeed,
			Description: "Add RSS feed to watch",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        FeedName,
					Description: "RSS feed name",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        FeedURL,
					Description: "RSS feed URL",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        TargetChannel,
					Description: "Channel to post updates in",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        UpdateInterval,
					Description: "How often to check for new items",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
			},
		},
		{
			Name:        ListRSSFeeds,
			Description: "List configured RSS feeds",
		},
		{
			Name:        SelectRSSFeedForDeletion,
			Description: "Remove RSS feed",
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
