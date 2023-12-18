package bot

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// registerOrUpdateServer checks if a guild is already registered in the database. If not,
// it creates it with sensibile defaults
func (bot *ModeratorBot) registerOrUpdateServer(g *discordgo.Guild, delete bool) error {
	status := sql.NullBool{Valid: true, Bool: true}
	if delete {
		status = sql.NullBool{Valid: true, Bool: false}
	}

	var registration ServerRegistration
	bot.DB.Find(&registration, g.ID)
	// The server registration does not exist, so we will create with defaults
	if registration.Name == "" {
		log.Infof("creating registration for new server: %s(%s)", g.Name, g.ID)
		registration = ServerRegistration{
			DiscordId: g.ID,
			Name:      g.Name,
			Active:    sql.NullBool{Valid: true, Bool: true},
			JoinedAt:  g.JoinedAt,
			Config: ServerConfig{
				Name:      g.Name,
				DiscordId: g.ID,
			},
		}
		tx := bot.DB.Model(&ServerRegistration{}).Where(&ServerRegistration{DiscordId: g.ID}).Save(&registration)

		// We only expect one server to be updated at a time. Otherwise, return an error
		if tx.RowsAffected != 1 {
			return fmt.Errorf("did not expect %v rows to be affected updating "+
				"server registration for server: %v(%v)", fmt.Sprintf("%v", tx.RowsAffected), g.Name, g.ID)
		}
	}

	// Update the registration if the DB is wrong or if the server
	// was deleted (the bot left) or if JoinedAt is not set
	// (field was added later so early registrations won't have it)
	if registration.Active != status || registration.JoinedAt.IsZero() {
		log.Debugf("updating server %s", g.Name)
		bot.DB.Model(&ServerRegistration{}).
			Where(&ServerRegistration{DiscordId: registration.DiscordId}).
			Updates(&ServerRegistration{Active: status, JoinedAt: g.JoinedAt})
		_ = bot.updateServersWatched()
	}

	// Called with no arguments, this only updates registration
	// for global commands.
	bot.UpdateCommands()

	return nil
}

// UpdateInactiveRegistrations goes through every server registration and
// updates the DB as to whether or not it's active
func (bot *ModeratorBot) UpdateInactiveRegistrations(activeGuilds []*discordgo.Guild) {
	var sr []ServerRegistration
	var inactiveRegistrations []string
	bot.DB.Find(&sr)

	// Check all registrations for whether or not the server is active
	for _, reg := range sr {
		active := false
		for _, g := range activeGuilds {
			if g.ID == reg.DiscordId {
				active = true
				break
			}
		}

		// If the server isn't found in activeGuilds, then we have a config
		// for a server we're not in anymore
		if !active {
			inactiveRegistrations = append(inactiveRegistrations, reg.DiscordId)
		}
	}

	// Since active servers will set Active to true, we will
	// set the rest of the servers as inactive.
	registrations := int64(len(inactiveRegistrations))
	if registrations > 0 {
		tx := bot.DB.Model(&ServerRegistration{}).Where(inactiveRegistrations).
			Updates(&ServerRegistration{Active: sql.NullBool{Valid: true, Bool: false}})

		if tx.RowsAffected != registrations {
			log.Errorf("unexpected number of rows affected updating %v inactive "+
				"server registrations, rows updated: %v",
				registrations, tx.RowsAffected)
		}
	}
}

// updateServersWatched updates the servers watched value
// in both the local bot stats and in the database. It is allowed to fail
func (bot *ModeratorBot) updateServersWatched() error {
	var serversConfigured, serversActive int64
	bot.DB.Model(&ServerRegistration{}).Where(&ServerRegistration{}).Count(&serversConfigured)
	serversActive = int64(len(bot.DG.State.Ready.Guilds))
	log.Debugf("total number of servers configured: %v, connected servers: %v", serversConfigured, serversActive)

	updateStatusData := &discordgo.UpdateStatusData{Status: "online"}
	updateStatusData.Activities = make([]*discordgo.Activity, 1)
	updateStatusData.Activities[0] = &discordgo.Activity{
		Name: fmt.Sprintf("%v %v", serversActive, handlePlural("server", "s", int(serversActive))),
		Type: discordgo.ActivityTypeWatching,
		URL:  moderaterRepoUrl,
	}

	log.Debug("updating discord bot status")
	err := bot.DG.UpdateStatusComplex(*updateStatusData)
	if err != nil {
		return fmt.Errorf("unable to update discord bot status: %w", err)
	}

	return nil
}
