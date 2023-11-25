package bot

import "fmt"

func (bot *ModeratorBot) GetUserReputation(userID string) (reputation int64, err error) {
	user := ModeratedUser{}
	tx := bot.DB.Model(&ModeratedUser{}).
		Where(&ModeratedUser{UserID: userID}).
		First(&user)

	if tx.RowsAffected > 1 {
		return 0, fmt.Errorf("unexpected number of rows affected getting user reputation: %v", tx.RowsAffected)
	}
	return user.Reputation, nil
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
