package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Takes user-submitted notes and adds it to an in-progress report
func (bot *ModeratorBot) SaveEvidenceNotes(i *discordgo.InteractionCreate) {
	notes := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	// If notes exist, update. If not, add new
	for idx, field := range i.Interaction.Message.Embeds[0].Fields {
		if field.Name == globals.Notes {
			i.Interaction.Message.Embeds[0].Fields[idx].Value = notes
			break
		}

		if idx == len(i.Interaction.Message.Embeds[0].Fields)-1 {
			i.Interaction.Message.Embeds[0].Fields = append(i.Interaction.Message.Embeds[0].Fields, &discordgo.MessageEmbedField{
				Name:  globals.Notes,
				Value: notes,
			})
		}
	}

	_ = bot.DG.InteractionRespond(i.Interaction,
		bot.DocumentBehaviorFromMessage(i, i.Interaction.Message))
}
