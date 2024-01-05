package bot

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const moderatorRepoUrl string = "https://github.com/tyzbit/go-discord-modtools"

// typeInChannel sets the typing indicator for a channel. The indicator is cleared
// when a message is sent or after some number of seconds.
func (bot *ModeratorBot) typeInChannel(channel chan bool, channelID string) {
	for {
		select {
		case <-channel:
			return
		default:
			if err := bot.DG.ChannelTyping(channelID); err != nil {
				log.Error("unable to set typing indicator: ", err)
			}
			time.Sleep(time.Second * 5)
		}
	}
}

// isAllowed returns a boolean if the user is in the preselected group
// that should have access to the bot
// PLANNED: or if the user is a server owner.
func (bot *ModeratorBot) isAllowed(cfg GuildConfig, member *discordgo.Member) bool {
	// Allow if no role has been set
	if cfg.ModeratorRoleSettingID == "" {
		log.Infof("Allowing %s(%s) to use function because moderator role is not defined in server %s(%s)",
			member.User.Username,
			member.User.ID,
			cfg.Name,
			cfg.ID,
		)
		return true
	}

	for _, roleID := range member.Roles {
		if roleID == cfg.ModeratorRoleSettingID {
			return true
		}
	}
	return false
}

// updateServersWatched updates the servers watched value
// in both the local bot stats and in the database. It is allowed to fail
func (bot *ModeratorBot) updateServersWatched() error {
	var GuildsConfigured, GuildsActive int64
	bot.DB.Where(&GuildConfig{}).Count(&GuildsConfigured)
	bot.DB.Where(&GuildConfig{Active: sql.NullBool{Bool: true, Valid: true}}).Count(&GuildsActive)
	log.Debugf("total number of servers configured: %v, connected servers: %v", GuildsConfigured, GuildsActive)

	updateStatusData := &discordgo.UpdateStatusData{Status: "online"}
	updateStatusData.Activities = make([]*discordgo.Activity, 1)
	updateStatusData.Activities[0] = &discordgo.Activity{
		Name: fmt.Sprintf("%v %v", GuildsActive, handlePlural("server", "s", int(GuildsActive))),
		Type: discordgo.ActivityTypeWatching,
		URL:  moderatorRepoUrl,
	}

	log.Debug("updating discord bot status")
	err := bot.DG.UpdateStatusComplex(*updateStatusData)
	if err != nil {
		return fmt.Errorf("unable to update discord bot status: %w", err)
	}

	return nil
}
