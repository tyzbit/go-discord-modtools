package bot

import (
	"fmt"

	"github.com/google/uuid"
)

func (bot *ModeratorBot) GetModeratedUser(serverID string, userID string) (moderatedUser ModeratedUser) {
	_ = bot.DB.Model(&ModeratedUser{}).
		Where(&ModeratedUser{UserID: userID}).
		First(&moderatedUser)

	guild, _ := bot.DG.Guild(serverID)
	user, _ := bot.DG.User(userID)
	// Create the moderated user if they do not exist
	if moderatedUser.UserID == "" {
		moderatedUser = ModeratedUser{
			UUID:       uuid.New().String(),
			UserName:   user.Username,
			UserID:     userID,
			ServerID:   serverID,
			ServerName: guild.Name,
			Reputation: 0,
		}
		bot.UpdateModeratedUser(moderatedUser)
	}
	return moderatedUser
}

// UpdateModeratedUser updates moderated user status in the database.
// It is allowed to fail
func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Updates(&u)

	if tx.RowsAffected > 1 {
		return fmt.Errorf("did not update one user row as expected, "+
			"updated %v rows for user %s(%s)",
			tx.RowsAffected, u.UserName, u.UserID)
	} else if tx.RowsAffected == 0 {
		// This user doesn't exist, so create them
		tx := bot.DB.Model(&ModeratedUser{}).Where(&ModeratedUser{UserID: u.UserID}).Create(&u)
		if tx.RowsAffected != 1 {
			return fmt.Errorf("did not create one user row as expected, "+
				"updated %v rows for user %s(%s)",
				tx.RowsAffected, u.UserName, u.UserID)
		}
	}
	return nil
}
