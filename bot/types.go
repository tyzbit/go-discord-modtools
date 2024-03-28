package bot

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// TODO: Change everything to Guild[...] to match DiscordGo
// TODO: Use pointers for nullable fields

// Events

// Handlers
// ModeratorBot is the main type passed around throughout the code
// It has many functions for overall bot management
type ModeratorBot struct {
	DB     *gorm.DB
	DG     *discordgo.Session
	Config ModeratorBotConfig
}

// ModeratorBotConfig is attached to ModeratorBot so config settings can be
// accessed easily
type ModeratorBotConfig struct {
	DBType     string `env:"DB_TYPE"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBUser     string `env:"DB_USER"`
	LogLevel   string `env:"LOG_LEVEL"`
	Token      string `env:"TOKEN"`
}

// Servers are called Guilds in the Discord API
type GuildConfig struct {
	ID                       string          `pretty:"Server ID" gorm:"type:varchar(191)"`
	Name                     string          `pretty:"Server Name" gorm:"default:default"`
	Active                   *bool           `pretty:"Bot is active in the server" gorm:"default:true"`
	EvidenceChannelSettingID string          `pretty:"Channel to document evidence in"`
	ModeratorRoleSettingID   string          `pretty:"Role for moderators"`
	CustomCommands           []CustomCommand `pretty:"Custom commands"`
	Polls                    []Poll          `pretty:"Active polls"`
}

// Custom commands registered with a specific server
// ID is registered with the Discord API
// GuildConfigID is the server ID where the command is registered
type CustomCommand struct {
	ID            string
	GuildConfigID string `gorm:"type:varchar(191)"`
	Name          string
	Description   string
	Content       string
}

// This is everything needed to track polls
type Poll struct {
	// This matches the message ID from discord when the poll was first posted
	ID            string `gorm:"type:varchar(191)"`
	GuildConfigID string
	Votes         []Vote
}

// Vote is one of 0 or more votes on a Poll
type Vote struct {
	ID            uint
	PollID        string
	GuildConfigID string
	// This is the ID of the button for the choice, used in a few places
	// It is only unique per PollID, essentially just the number of the poll choice
	CustomID string
	UserID   string
}

// A ModeratedUser represents a specific user/server combination
// that serves to record events and a "Reputation" which is only
// visible to people who are in the moderator's configurable role.
type ModeratedUser struct {
	ID               string
	GuildId          string
	GuildName        string
	UserID           string
	UserName         string
	ModerationEvents []ModerationEvent
	Reputation       *int64 `gorm:"default:1"`
}

// This is the representation of a moderation action
type ModerationEvent struct {
	ID                 uint `gorm:"primaryKey"`
	ModeratedUserID    string
	GuildId            string
	GuildName          string
	UserID             string
	UserName           string
	Notes              string
	PreviousReputation *int64
	CurrentReputation  *int64
	ModeratorID        string
	ModeratorName      string
	ReportURL          string
}

// Stats
// botStats is read by structToPrettyDiscordFields and converted
// into a slice of *discordgo.MessageEmbedField
type botStats struct {
	ModerateRequests int64 `pretty:"Times the bot has been called"`
	MessagesSent     int64 `pretty:"Messages Sent"`
	Interactions     int64 `pretty:"Interactions with the bot"`
	GuildsActive     int64 `pretty:"Active servers"`
	GuildsConfigured int64 `pretty:"Configured servers" global:"true"`
}
