package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (bot *ModeratorBot) isAllowed(sc ServerConfig, member *discordgo.Member) bool {
	// Allow if no role has been set
	if sc.ModeratorRoleSettingID == "" {
		log.Infof("Allowing %s(%s) to use function because moderator role is not defined in server %s(%s)",
			member.User.Username,
			member.User.ID,
			sc.Name,
			sc.DiscordId,
		)
		return true
	}

	for _, roleID := range member.Roles {
		if roleID == sc.ModeratorRoleSettingID {
			return true
		}
	}
	return false
}
