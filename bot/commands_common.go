package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tyzbit/go-discord-modtools/globals"
)

func (bot *ModeratorBot) DocumentBehaviorFromMessage(i *discordgo.InteractionCreate, message *discordgo.Message) (resp *discordgo.InteractionResponse) {
	user := bot.GetModeratedUser(i.GuildID, message.Author.ID)
	var fields []*discordgo.MessageEmbedField
	if len(message.Embeds) > 0 {
		fields = message.Embeds[0].Fields
	} else {
		fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Original message timestamp",
				Value:  fmt.Sprintf("%s (<t:%v:R>)", message.Timestamp.Format(time.RFC1123Z), message.Timestamp.Unix()),
				Inline: false,
			},
			{
				Name:   "Author of message",
				Value:  fmt.Sprintf("<@%s> (Reputation: %v)", user.UserID, user.Reputation.Int64),
				Inline: true,
			},
			{
				Name:   "Link to original message",
				Value:  fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.Interaction.GuildID, message.ChannelID, message.ID),
				Inline: true,
			},
			{
				Name: "Collected by",
				Value: fmt.Sprintf(`<@%s>
									%s (<t:%v:R>)`,
					i.Interaction.Member.User.ID,
					time.Now().Format(time.RFC1123Z),
					time.Now().Unix()),
			},
			{
				Name:  globals.OriginalMessageContent,
				Value: message.Content,
			},
		}

		if len(message.Attachments) > 0 {
			attachmentList := ``
			for _, attachment := range message.Attachments {
				attachmentList = attachmentList + fmt.Sprintf("%s: %v\n", attachment.Filename, attachment.URL)
			}
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf(globals.Attachments+" (%v)", len(message.Attachments)),
				Value: attachmentList,
			})
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Create evidence report",
					Description: fmt.Sprintf("Document user behavior for for <@%v> - good, bad, or noteworthy", message.Author.ID),
					Color:       globals.Purple,
					Fields:      fields,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: globals.DecreaseUserReputation,
						Label:    globals.DecreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.IncreaseUserReputation,
						Label:    globals.IncreaseUserReputation,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.TakeEvidenceNotes,
						Label:    globals.TakeEvidenceNotes,
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						CustomID: globals.SubmitReport,
						Label:    globals.SubmitReport,
						Style:    discordgo.PrimaryButton,
					},
				},
			}},
		},
	}
}
