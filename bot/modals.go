package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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

	err := bot.DG.InteractionRespond(i.Interaction,
		bot.DocumentBehaviorFromMessage(i, i.Interaction.Message))
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Takes user-submitted command and adds it to the server config and registers
// it in the guild
func (bot *ModeratorBot) SaveCustomSlashCommand(i *discordgo.InteractionCreate) {
	name := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	description := i.Interaction.ModalSubmitData().
		Components[1].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	content := i.Interaction.ModalSubmitData().
		Components[2].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	customCommand := CustomCommand{
		DiscordId:   i.GuildID,
		Name:        name,
		Description: description,
		Content:     content,
	}

	sc := bot.getServerConfig(i.GuildID)
	sc.CustomCommands = append(sc.CustomCommands, customCommand)
	tx := bot.DB.Save(&sc)
	if tx.RowsAffected != 1 {
		log.Warn("something other than one row affected when updating custom slash command")
	}

	bot.RegisterCustomCommandHandler(sc)

	// Update registered commands
	bot.UpdateCommands()

	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: fmt.Sprintf("Custom command /%s created", customCommand.Name),
					Description: fmt.Sprintf(`Description:
					%s
					
					Content:
						%s`, customCommand.Description, customCommand.Content),
					Color: globals.Purple,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("error responding to custom slash command creation, err: %v", err)
	}
}
