package bot

import (
	"database/sql"
)

// Returns a ModeratedUser record from the DB using server and user ID
// (a user can be in multiple servers)
func (bot *ModeratorBot) GetModeratedUser(GuildId string, userID string) (moderatedUser ModeratedUser) {
	guild, _ := bot.DG.Guild(GuildId)
	user, _ := bot.DG.User(userID)
	moderatedUser = ModeratedUser{
		UserName:   user.Username,
		UserID:     userID,
		GuildId:    GuildId,
		ID:         GuildId + userID,
		GuildName:  guild.Name,
		Reputation: sql.NullInt64{Int64: 1, Valid: true},
	}
	bot.DB.Where(&ModeratedUser{ID: GuildId + userID}).FirstOrCreate(&moderatedUser)
	return moderatedUser
}
