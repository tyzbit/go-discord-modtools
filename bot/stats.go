package bot

import (
	"database/sql"
)

// getGlobalStats calls the database to get global stats for the bot
// The output here is not appropriate to send to individual servers, except
// for ServersActive
func (bot *ModeratorBot) getGlobalStats() botStats {
	var ModerateRequests, MessagesSent, Interactions, URLsModerated, ServersConfigured, ServersActive int64

	bot.DB.Model(&ModerationEvent{}).Count(&URLsModerated)
	bot.DB.Model(&ServerRegistration{}).Count(&ServersConfigured)
	bot.DB.Find(&ServerRegistration{}).Where(&ServerRegistration{
		Active: sql.NullBool{Valid: true, Bool: true}}).Count(&ServersActive)

	return botStats{
		ModerateRequests:  ModerateRequests,
		MessagesSent:      MessagesSent,
		Interactions:      Interactions,
		ServersConfigured: ServersConfigured,
		ServersActive:     ServersActive,
	}
}

// getServerStats gets the stats for a particular server with ID serverId
// If you want global stats, use getGlobalStats()
func (bot *ModeratorBot) getServerStats(serverId string) botStats {
	var ModerateRequests, MessagesSent, Interactions, ServersActive int64

	bot.DB.Model(&ModerationEvent{}).Where(&ModerationEvent{ServerID: serverId}).Count(&ModerateRequests)
	bot.DB.Model(&ServerRegistration{}).Where(&ServerRegistration{}).Count(&ServersActive)

	return botStats{
		ModerateRequests: ModerateRequests,
		MessagesSent:     MessagesSent,
		Interactions:     Interactions,
		ServersActive:    ServersActive,
	}
}
