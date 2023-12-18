package bot

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Events

// This is the representation of a moderation action
type ModerationEvent struct {
	CreatedAt          time.Time
	UUID               string `gorm:"primaryKey;uniqueIndex"`
	ServerID           string `gorm:"index"`
	ServerName         string
	UserID             string
	UserName           string
	Notes              string
	PreviousReputation sql.NullInt64
	CurrentReputation  sql.NullInt64
	ModeratorID        string
	ModeratorName      string
	ReportURL          string
}

// A ModeratedUser represents a specific user/server combination
// that serves to track changes
type ModeratedUser struct {
	CreatedAt        time.Time
	UUID             string `gorm:"primaryKey;uniqueIndex"`
	ServerID         string `gorm:"index"`
	ServerName       string
	UserID           string
	UserName         string
	ModerationEvents []ModerationEvent `gorm:"foreignKey:UUID"`
	Reputation       sql.NullInt64     `gorm:"default:1"`
}

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
	AdminIds   []string `env:"ADMINISTRATOR_IDS"`
	DBHost     string   `env:"DB_HOST"`
	DBName     string   `env:"DB_NAME"`
	DBPassword string   `env:"DB_PASSWORD"`
	DBUser     string   `env:"DB_USER"`
	LogLevel   string   `env:"LOG_LEVEL"`
	Token      string   `env:"TOKEN"`
}

// Servers
type ServerRegistration struct {
	gorm.Model
	DiscordId string
	Name      string
	JoinedAt  time.Time
	Active    sql.NullBool `pretty:"Bot is active in the server" gorm:"default:true"`
	Config    ServerConfig
}

// Configuration for each server, changed with `/settings`
type ServerConfig struct {
	gorm.Model
	ServerRegistrationID     uint
	DiscordId                string          `pretty:"ServerID"`
	Name                     string          `pretty:"Server Name" gorm:"default:default"`
	EvidenceChannelSettingID string          `pretty:"Channel to document evidence in"`
	ModeratorRoleSettingID   string          `pretty:"Role for moderators"`
	CustomCommands           []CustomCommand `pretty:"Custom commands"`
}

// Custom commands registered with a specific server
type CustomCommand struct {
	gorm.Model
	ServerConfigID uint
	DiscordId      string
	Name           string
	Description    string
	Content        string
}

// Stats
// botStats is read by structToPrettyDiscordFields and converted
// into a slice of *discordgo.MessageEmbedField
type botStats struct {
	ModerateRequests  int64 `pretty:"Times the bot has been called"`
	MessagesSent      int64 `pretty:"Messages Sent"`
	Interactions      int64 `pretty:"Interactions with the bot"`
	ServersActive     int64 `pretty:"Active servers"`
	ServersConfigured int64 `pretty:"Configured servers" global:"true"`
}
