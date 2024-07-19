package bot

const (
	// Commands (max 32 in length)
	// These all function as IDs, they are sometimes shown to the user
	// They must be unique among similar types (ex: all command IDs must be unique)

	// TODO: make this i8n compatible (one task of which is to ensure changing the
	// translated English words does not break the command IDs that are registered
	// from them

	// Premade Option IDs (semi-reusable)
	// TODO: actions that delete messages, ban users etc
	UserOption    = "user"
	ChannelOption = "channel"
	MessageOption = "message"
	ReasonOption  = "reason"

	// Constants
	MaxCommandContentLength     = 32   // https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
	MaxMessageContentLength     = 2000 // https://discord.com/developers/docs/resources/channel#create-message
	MaxDescriptionContentLength = 100  // https://discord.com/developers/docs/interactions/application-commands#application-command-object

	// Colors
	FrenchGray = 13424349
	Purple     = 7283691
	DarkRed    = 9109504
	Green      = 4306266

	// If images have this in front of their name, they're spoilered
	Spoiler = "SPOILER_"

	// URLs
	MessageURLTemplate = "https://discord.com/channels/%s/%s/%s"
)
