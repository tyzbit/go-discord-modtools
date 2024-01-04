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
		GuildID: g.ID,
		Name:    g.Name,
		Active:  sql.NullBool{Valid: true, Bool: true},
	}
	tx := bot.DB.Where(&GuildConfig{GuildID: g.ID}).FirstOrCreate(&reg)
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
			if g.ID == reg.GuildID {
				active = true
				break
			}
		}

		// If the server isn't found in activeGuilds, then we have a config
		// for a server we're not in anymore
		if !active {
			inactiveRegistrations = append(inactiveRegistrations, reg.GuildID)
		}
	}

	// Since active servers will set Active to true, we will
	// set the rest of the servers as inactive.
	registrations := int64(len(inactiveRegistrations))
	if registrations > 0 && len(inactiveRegistrations) > 0 {
		tx := bot.DB.Model(&GuildConfig{}).Where(inactiveRegistrations)
		tx.Updates(&GuildConfig{Active: sql.NullBool{Valid: true, Bool: false}})

		if tx.RowsAffected != registrations {
			log.Errorf("unexpected number of rows affected updating %v inactive "+
				"server registrations, rows updated: %v",
				registrations, tx.RowsAffected)
		}
	}
}
