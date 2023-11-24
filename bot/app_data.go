package bot

import "fmt"

// UpdateModeratedUser updates moderated user status in the database.
// It is allowed to fail
func (bot *ModeratorBot) UpdateModeratedUser(u ModeratedUser) error {
	tx := bot.DB.Model(&ModeratedUser{}).
		Where(&ModeratedUser{UserID: u.UserID}).Updates(u)

	if tx.RowsAffected != 1 {
		return fmt.Errorf("did not update one user row as expected, "+
			"updated %v rows for user %s(%s)",
			tx.RowsAffected, u.UserName, u.UserID)
	}
	return nil
}
