package globals

import (
	"github.com/bwmarrin/discordgo"
)

const (
	// Commands (max 32 in length)
	// These all function as IDs, they are sometimes shown to the user
	// They must be unique among similar types (ex: all command IDs must be unique)

	// Chat commands
	Help                              = "help"
	Stats                             = "stats"
	Settings                          = "settings"
	GetUserInfoFromChatCommandContext = "query"

	// User commands
	GetUserInfoFromUserContext      = "Check info"
	DocumentBehaviorFromUserContext = "Save evidence"

	// Message commands
	GetUserInfoFromMessageContext      = "Get user info"
	DocumentBehaviorFromMessageContext = "Save as evidence"

	// Premade Option IDs (semi-reusable)
	// TODO: actions that delete messages, ban users etc
	// TODO: (extra credit) make a command that manages custom commands that drop a premade message
	UserOption    = "user"
	ChannelOption = "channel"
	MessageOption = "message"
	ReasonOption  = "reason"

	// Buttons
	// Settings (the names affect the column names in the DB)
	EvidenceChannelSettingID = "Evidence channel"
	ModeratorRoleSettingID   = "Moderator role"

	// Moderation buttons
	IncreaseUserReputation = "⬆️ Reputation"
	DecreaseUserReputation = "⬇️ Reputation"
	TakeEvidenceNotes      = "Add notes"
	SubmitReport           = "Submit report"

	// Modals
	SaveEvidenceNotes = "Save evidence notes"

	// Modal options
	EvidenceNotes = "Evidence notes"

	// Colors
	FrenchGray = 13424349
	Purple     = 7283691
	DarkRed    = 9109504

	// Text fragments
	CurrentReputation      = "Current reputation"
	PreviousReputation     = "Previous reputation"
	ModerationSuccessful   = "Moderation action saved"
	ModerationUnSuccessful = "There was a problem saving moderation action"
	OriginalMessageContent = "Content of original message"
	Attachments            = "Attachments"
	Notes                  = "Notes"
	// If images have this in front of their name, they're spoilered
	Spoiler = "SPOILER_"

	// URLs
	MessageURLTemplate = "https://discord.com/channels/%s/%s/%s"

	// Shown to the user when `/help` is called
	BotHelpText = `**Usage**
	Right click (Desktop) or long press (mobile) a message and select Apps, then select a moderation action for the message.

You can also right-click (or long press) a message and use "Get snapshot" to get a message with snapshots that only you can see- You can also right-click (or long press) a message and use "Get snapshot" to get a message with snapshots that only you can see or select "Take snapshot" to take a snapshot of the live page.

Configure the bot:

` + "`/settings`" + `

Look up a user (more coming soon):

` + "`/query`" + `

Directly moderate up a user (more coming soon):

` + "`/moderate`" + `

Get stats for the bot:

` + "`/stats`" + `

Get this help message:

` + "`/help`"
)

var (
	// Enabled takes a boolean and returns "enabled" or "disabled"
	Enabled = map[bool]string{
		true:  "enabled",
		false: "disabled",
	}
	// Button style takes a boolean and returns a colorized button if true
	ButtonStyle = map[bool]discordgo.ButtonStyle{
		true:  discordgo.PrimaryButton,
		false: discordgo.SecondaryButton,
	}
	SettingFailedResponseMessage = "Error changing setting"

	// These objects are used to register chat commands, so if it's
	// not in here, it won't get registered properly.
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
			Name:        Stats,
			Description: "Show bot stats",
		},
		{
			Name:        Settings,
			Description: "Change settings",
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

	RegisteredCommands = make([]*discordgo.ApplicationCommand,
		len(ChatCommands)+len(UserCommands)+len(MessageCommands))

	// This object is used to register chat commands, so if it's
	// not in here, it won't get registered properly.
	Commands = append(append(ChatCommands, UserCommands...), MessageCommands...)
)
