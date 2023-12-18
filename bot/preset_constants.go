package bot

const (
	// Commands (max 32 in length)
	// These all function as IDs, they are sometimes shown to the user
	// They must be unique among similar types (ex: all command IDs must be unique)

	// TODO: make this i8n compatible (one task of which is to ensure changing the
	// translated English words does not break the command IDs that are registered
	// from them

	// Chat commands
	Help                              = "help"
	Stats                             = "stats"
	Settings                          = "settings"
	GetUserInfoFromChatCommandContext = "query"
	AddCustomSlashCommand             = "addcommand"

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
	IncreaseUserReputation      = "⬆️ Reputation"
	DecreaseUserReputation      = "⬇️ Reputation"
	ShowEvidenceCollectionModal = "Add notes"
	SubmitReport                = "Submit report"

	// Modals
	SaveEvidenceNotes      = "Save evidence notes"
	SaveCustomSlashCommand = "Save custom slash command"

	// Modal options
	EvidenceNotes                 = "Evidence notes"
	CustomSlashName               = "Name for this custom slash command"
	CustomSlashCommandDescription = "Description for this custom slash command"
	CustomSlashCommandContent     = "Message to paste if this command is used"
	CustomCommandIdentifier       = "Custom command: "

	// Colors
	FrenchGray = 13424349
	Purple     = 7283691
	DarkRed    = 9109504

	// Constants
	MaxCommandContentLength     = 32   // https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
	MaxMessageContentLength     = 2000 // https://discord.com/developers/docs/resources/channel#create-message
	MaxDescriptionContentLength = 100  // https://discord.com/developers/docs/interactions/application-commands#application-command-object

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

	// Errors
	SettingFailedResponseMessage = "Error changing setting"
	UnexpectedRowsAffected       = "unexpected number of rows affected getting user reputation: %v"

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
