package globals

import (
	"github.com/bwmarrin/discordgo"
)

const (
	// Commands (max 32 in length)
	Stats                                          = "stats"
	Query                                          = "query"
	Settings                                       = "settings"
	ModerateUser                                   = "Moderate user"
	ModerateMessage                                = "Moderate message"
	CollectMessageAsEvidence                       = "Collect evidence"
	CollectMessageAsEvidenceThenIncreaseReputation = "Increse user rep"
	CollectMessageAsEvidenceThenDecreaseReputation = "Decrease user rep"
	CheckUserReputationUsingMessage                = "Check user rep"
	CheckUserReputation                            = "Check reputation"
	// TODO: actions that delete messages, ban users etc
	// TODO: (extra credit) make a command that manages custom commands that drop a premade message
	Help = "help"

	UserOption    = "user"
	MessageOption = "message"
	ReasonOption  = "reason"

	// Settings
	NotifyWhenReputationIsBelowSetting = "Low rep notification"
	NotifyWhenReputationIsAboveSetting = "High rep notification"
	EvidenceChannelSettingID           = "Evidence channel"

	// Modal
	ModerateModal       = "Moderate user modal"
	ModerateModalReason = "Reason"

	// Colors
	FrenchGray = 13424349
	Purple     = 7283691

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
	Commands                     = []*discordgo.ApplicationCommand{
		{
			Name:        Help,
			Description: "How to use this bot",
		},
		{
			Name:        Query,
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
		{
			Name: CollectMessageAsEvidence,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: CollectMessageAsEvidenceThenIncreaseReputation,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: CollectMessageAsEvidenceThenDecreaseReputation,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: CheckUserReputationUsingMessage,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: CheckUserReputation,
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: ModerateMessage,
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: ModerateUser,
			Type: discordgo.UserApplicationCommand,
		},
	}
	RegisteredCommands = make([]*discordgo.ApplicationCommand, len(Commands))
)
