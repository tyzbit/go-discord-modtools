package bot

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// registerOrUpdateServer checks if a guild is already registered in the database. If not,
// it creates it with sensibile defaults
func (bot *ModeratorBot) registerOrUpdateServer(g *discordgo.Guild, delete bool) error {
	reg := GuildConfig{
		ID:     g.ID,
		Name:   g.Name,
		Active: sql.NullBool{Valid: true, Bool: true},
	}
	tx := bot.DB.Where(&GuildConfig{ID: g.ID}).FirstOrCreate(&reg)
	tx.Commit()

	// Called with no arguments, this only updates registration
	// for global commands.
	err := bot.UpdateCommands()
	if err != nil {
		return err
	}

	return nil
}

// UpdateInactiveRegistrations goes through every server registration and
// updates the DB as to whether or not it's active
func (bot *ModeratorBot) UpdateInactiveRegistrations(activeGuilds []*discordgo.Guild) {
	var sr []GuildConfig
	var inactiveRegistrations []string
	bot.DB.Find(&sr)

	// Check all registrations for whether or not the server is active
	for _, reg := range sr {
		active := false
		for _, g := range activeGuilds {
			if g.ID == reg.ID {
				active = true
				break
			}
		}

		// If the server isn't found in activeGuilds, then we have a config
		// for a server we're not in anymore
		if !active {
			inactiveRegistrations = append(inactiveRegistrations, reg.ID)
		}
	}

	// Since active servers will set Active to true, we will
	// set the rest of the servers as inactive.
	if len(inactiveRegistrations) > 0 {
		tx := bot.DB.Model(&GuildConfig{}).Where("id IN ?", inactiveRegistrations).
			Updates(map[string]interface{}{"active": sql.NullBool{Valid: true, Bool: false}})

		if tx.RowsAffected != int64(len(inactiveRegistrations)) {
			log.Errorf("unexpected number of rows affected updating %v inactive "+
				"server registrations, rows updated: %v",
				inactiveRegistrations, tx.RowsAffected)
		}
	}
}
