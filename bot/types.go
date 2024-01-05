package bot

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// TODO: Change everything to Guild[...] to match DiscordGo

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
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBUser     string `env:"DB_USER"`
	LogLevel   string `env:"LOG_LEVEL"`
	Token      string `env:"TOKEN"`
}

// Servers
// TODO: Combine registration and config
type GuildConfig struct {
	CreatedAt                time.Time
	UpdatedAt                time.Time
	GuildID                  string          `pretty:"Server ID" gorm:"primarykey"`
	Name                     string          `pretty:"Server Name" gorm:"default:default"`
	Active                   sql.NullBool    `pretty:"Bot is active in the server" gorm:"default:true"`
	EvidenceChannelSettingID string          `pretty:"Channel to document evidence in"`
	ModeratorRoleSettingID   string          `pretty:"Role for moderators"`
	CustomCommands           []CustomCommand `pretty:"Custom commands"`
}

// Custom commands registered with a specific server
type CustomCommand struct {
	gorm.Model
	GuildConfigID string
	DiscordId     string
	Name          string
	Description   string
	Content       string
}

// A ModeratedUser represents a specific user/server combination
// that serves to track changes
type ModeratedUser struct {
	ID               string `gorm:"unique"`
	GuildId          string
	GuildName        string
	UserID           string
	UserName         string
	ModerationEvents []ModerationEvent
	Reputation       sql.NullInt64 `gorm:"default:1"`
}

// This is the representation of a moderation action
type ModerationEvent struct {
	ModeratedUserID    string
	GuildId            string
	GuildName          string
	UserID             string
	UserName           string
	Notes              string
	PreviousReputation sql.NullInt64
	CurrentReputation  sql.NullInt64
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
