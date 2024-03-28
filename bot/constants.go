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
	RemoveCustomSlashCommand          = "deletecommand"
	CreatePoll                        = "poll"

	// User commands
	GetUserInfoFromUserContext      = "Check info"
	DocumentBehaviorFromUserContext = "Save evidence"

	// Message commands
	GetUserInfoFromMessageContext      = "Get user info"
	DocumentBehaviorFromMessageContext = "Save as evidence"

	// Premade Option IDs (semi-reusable)
	// TODO: actions that delete messages, ban users etc
	UserOption          = "user"
	ChannelOption       = "channel"
	MessageOption       = "message"
	ReasonOption        = "reason"
	NumberOfPollChoices = "choices"

	// Buttons
	// Settings (the names affect the column names in the DB)
	EvidenceChannelSettingID = "Evidence channel"
	ModeratorRoleSettingID   = "Moderator role"
	// Poll buttons
	EndPoll          = "End poll"
	PollOptionPrefix = "Poll option "

	// Moderation buttons
	IncreaseUserReputation      = "⬆️ Reputation"
	DecreaseUserReputation      = "⬇️ Reputation"
	ShowEvidenceCollectionModal = "Add notes"
	SubmitReport                = "Submit report"

	// Modals
	SaveEvidenceNotes        = "Save evidence notes"
	SaveCustomSlashCommand   = "Save custom slash command"
	DeleteCustomSlashCommand = "Remove custom slash command"
	StartPoll                = "Start poll"

	// Modal options
	EvidenceNotes                 = "Evidence notes"
	CustomSlashName               = "Name for this custom slash command"
	CustomSlashCommandDescription = "Description for this custom slash command"
	CustomSlashCommandContent     = "Message to paste if this command is used"
	CustomCommandIdentifier       = "Custom command: "
	VoteEndTime                   = "Time when the vote should end"
	VoteOptions                   = "List of options for voting"
	PollName                      = "Name for the poll"

	// Colors
	FrenchGray = 13424349
	Purple     = 7283691
	DarkRed    = 9109504
	Green      = 4306266

	// Emoji
	StopEmoji = "🛑"

	// Constants
	MaxCommandContentLength     = 32   // https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
	MaxMessageContentLength     = 2000 // https://discord.com/developers/docs/resources/channel#create-message
	MaxDescriptionContentLength = 100  // https://discord.com/developers/docs/interactions/application-commands#application-command-object
	DefaultNumberOfPollOptions  = "3"

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
	In the Discord app, right click (Desktop) or long press (mobile) a message or user to see the available options.

Configure the bot (Highly recommended to do so only people in a specific role can use the moderation commands):

` + "`/settings`" + `

Look up a user (more coming soon):

` + "`/query`" + `

Add a custom command (simply posts your desired block of text, Markdown formatting enabled)

` + "`/addcommand`" + `

Remove a custom command

` + "`/deletecommand`" + `

Get this help message:

` + "`/help`"
)
