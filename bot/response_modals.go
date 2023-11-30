package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Takes user-submitted notes and adds it to an in-progress report
func (bot *ModeratorBot) SaveEvidenceNotes(i *discordgo.InteractionCreate) {
	notes := i.Interaction.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	i.Interaction.Message.Embeds[0].Fields = append(i.Interaction.Message.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:  "Notes",
		Value: notes,
	})

	_ = bot.DG.InteractionRespond(i.Interaction,
		bot.DocumentBehaviorFromMessage(i, i.Interaction.Message))
}
