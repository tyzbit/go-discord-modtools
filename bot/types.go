package bot

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// Events
// A MessageEvent is created when we receive a message
type MessageEvent struct {
	CreatedAt      time.Time
	UUID           string `gorm:"primaryKey" gorm:"uniqueIndex"`
	AuthorId       string `gorm:"index"`
	AuthorUsername string
	MessageId      string
	ChannelId      string
	ServerID       string `gorm:"index"`
	ServerName     string
}

// A InteractionEvent when a user interacts with an Embed
type InteractionEvent struct {
	CreatedAt        time.Time
	UUID             string `gorm:"primaryKey" gorm:"uniqueIndex"`
	UserID           string `gorm:"index"`
	Username         string
	InteractionId    string
	ChannelId        string
	ServerID         string `gorm:"index"`
	ServerName       string
	ModerationEvents []ModerationEvent `gorm:"foreignKey:UUID"`
}

// This is the representation of a moderation action
type ModerationEvent struct {
	CreatedAt  time.Time
	UUID       string `gorm:"primaryKey;uniqueIndex"`
	ServerID   string `gorm:"index"`
	ServerName string
	User       discordgo.User
	Message    discordgo.Message
	Action     string
	Reason     string
	Moderator  discordgo.User
}

type ModeratedUser struct {
	CreatedAt        time.Time
	UUID             string `gorm:"primaryKey;uniqueIndex"`
	ServerID         string `gorm:"index"`
	ServerName       string
	User             discordgo.User
	Message          discordgo.Message
	ModerationEvents []ModerationEvent `gorm:"foreignKey:UUID"`
	Reputation       int64
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
	DiscordId string `gorm:"primaryKey;uniqueIndex"`
	Name      string
	UpdatedAt time.Time
	JoinedAt  time.Time
	Active    sql.NullBool `pretty:"Bot is active in the server" gorm:"default:true"`
	Config    ServerConfig `gorm:"foreignKey:DiscordId"`
}

type ServerConfig struct {
	DiscordId                   string            `gorm:"primaryKey;uniqueIndex" pretty:"Server ID"`
	Name                        string            `pretty:"Server Name" gorm:"default:default"`
	NotifyWhenReputationIsBelow sql.NullInt32     `pretty:"Notify when a user's reputation falls below this" gorm:"default:5"`
	NotifyWhenReputationIsAbove sql.NullInt32     `pretty:"Notify when a user's reputation is greater than this" gorm:"default:3"`
	EvidenceChannel             discordgo.Channel `pretty:"Channel to document evidence in"`
	UpdatedAt                   time.Time
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
