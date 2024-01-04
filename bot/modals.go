package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Takes user-submitted notes and adds it to an in-progress report
func (bot *ModeratorBot) SaveEvidenceNotes(i *discordgo.InteractionCreate) {
	notes := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	// If notes exist, update. If not, add new
	for idx, field := range i.Interaction.Message.Embeds[0].Fields {
		if field.Name == Notes {
			i.Interaction.Message.Embeds[0].Fields[idx].Value = notes
			break
		}

		if idx == len(i.Interaction.Message.Embeds[0].Fields)-1 {
			i.Interaction.Message.Embeds[0].Fields = append(i.Interaction.Message.Embeds[0].Fields, &discordgo.MessageEmbedField{
				Name:  Notes,
				Value: notes,
			})
		}
	}

	err := bot.DG.InteractionRespond(i.Interaction,
		bot.GenerateEvidenceReportFromMessage(i, i.Interaction.Message))
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Takes user-submitted command and adds it to the server config and registers
// it in the guild
func (bot *ModeratorBot) SaveCustomSlashCommand(i *discordgo.InteractionCreate) {
	guild, _ := bot.DG.Guild(i.GuildID)
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
		GuildConfigID: guild.ID,
		Name:          strings.ToLower(name),
		Description:   CustomCommandIdentifier + description,
		Content:       content,
	}

	ird := discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf("Custom command /%s created", customCommand.Name),
				Description: fmt.Sprintf(`Description:
				%s
				
				Content:
					%s`, description, customCommand.Content),
				Color: Purple,
			},
		},
	}

	if len(customCommand.Description) > MaxDescriptionContentLength {
		ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Please limit the description to %v", MaxDescriptionContentLength-len(CustomCommandIdentifier)))
	} else if len(customCommand.Content) > MaxMessageContentLength {
		ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Please limit the message length to %v", MaxMessageContentLength))
	} else if strings.Contains(customCommand.Name, " ") {
		ird = *bot.generalErrorDisplayedToTheUser("Command names may not have spaces")
	} else {
		var cmds []CustomCommand
		bot.DB.Model(&CustomCommand{}).Where(&CustomCommand{GuildConfigID: guild.ID}).Create(&customCommand)

		// Register commands with Discord API
		bot.RegisterCustomCommandHandler(cmds)

		info := fmt.Sprintf(" from guild %s(%s)", guild.Name, guild.ID)
		log.Debugf("creating command '/%s'", name+info)
		_, err := bot.DG.ApplicationCommandCreate(bot.DG.State.User.ID,
			guild.ID,
			&discordgo.ApplicationCommand{
				Name:        customCommand.Name,
				Description: customCommand.Description,
				GuildID:     guild.ID,
			})
		if err != nil {
			ird = *bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Unable to create command, err: %v", err))
		}
	}

	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &ird,
	})
	if err != nil {
		log.Errorf("error responding to custom slash command creation, err: %v", err)
	}
}
