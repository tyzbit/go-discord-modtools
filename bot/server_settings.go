package bot

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const moderaterRepoUrl string = "https://github.com/tyzbit/go-discord-modtools"

// getServerConfig takes a guild ID and returns a ServerConfig object for that server
// If the config isn't found, it returns a default config
func (bot *ModeratorBot) getServerConfig(guildId string) ServerConfig {
	// Default server config in case guild lookup fails
	sc := ServerConfig{
		DiscordId:                "",
		Name:                     "",
		UpdatedAt:                time.Now(),
		EvidenceChannelSettingID: "",
		ModeratorRoleSettingID:   "",
	}
	// If this fails, we'll return a default server
	// config, which is expected
	bot.DB.Where(&ServerConfig{DiscordId: guildId}).Find(&sc)
	return sc
}

// updateServerSetting updates a server setting according to the
// column name (setting) and the value
func (bot *ModeratorBot) updateServerSetting(guildID string, setting string,
	value interface{}) (sc ServerConfig, success bool) {
	guild, err := bot.DG.Guild(guildID)
	if err != nil {
		log.Errorf("unable to look up server by id: %v", guildID)
		return sc, false
	}

	tx := bot.DB.Model(&ServerConfig{}).Where(&ServerConfig{DiscordId: guild.ID}).
		Update(setting, value)

	ok := true
	// We only expect one server to be updated at a time. Otherwise, return an error
	if tx.RowsAffected != 1 {
		log.Errorf("did not expect %v rows to be affected updating "+
			"server config for server: %v(%v)", fmt.Sprintf("%v", tx.RowsAffected), guild.Name, guild.ID)
		ok = false
	}
	return bot.getServerConfig(guildID), ok
}
