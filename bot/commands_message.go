package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Moderation modal menu
func (bot *ModeratorBot) DocumentBehaviorFromMessageContext(i *discordgo.InteractionCreate) {
	data := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if data.ID == "" {
		reason := "No message was provided"
		log.Warn(reason)
		_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bot.generalErrorDisplayedToTheUser(reason)})
		return
	}

	message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if data.ID == "" {
		reason := "No message was provided"
		log.Warn(reason)
		_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bot.generalErrorDisplayedToTheUser(reason)})
		return
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Original message timestamp",
			Value:  fmt.Sprintf("%s (<t:%v:R>)", message.Timestamp.Format(time.RFC1123Z), message.Timestamp.Unix()),
			Inline: false,
		},
		{
			Name:   "Author of message",
			Value:  fmt.Sprintf("<@%s>", message.Author.ID),
			Inline: true,
		},
		{
			Name:   "Link to original message",
			Value:  fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, message.ChannelID, message.ID),
			Inline: true,
		},
		{
			Name: "Collected by",
			Value: fmt.Sprintf("<@%s> at %s n(<t:%v:R>)",
				i.Interaction.Member.User.ID,
				time.Now().Format(time.RFC1123Z),
				time.Now().Unix()),
		},
		{
			Name:  "Content of original message",
			Value: message.Content,
		},
	}

	if len(message.Attachments) > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Attachments",
			Value: fmt.Sprintf("%v", len(message.Attachments)),
		})
		for _, attachment := range message.Attachments {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  attachment.Filename,
				Value: attachment.URL,
			})
		}
	}

	// TODO: more information

	// ms := discordgo.MessageSend{
	// 	Embeds: []*discordgo.MessageEmbed{{
	// 		Title: "Moder8s Evidence",
	// 		Description: fmt.Sprintf("Collected by <@%s> at %s \\n(<t:%v:R>)",
	// 			i.Interaction.Member.User.ID,
	// 			time.Now().Format(time.RFC1123Z),
	// 			time.Now().Unix(),
	// 		),
	// 		Fields: fields,
	// 		Color:  globals.Purple,
	// 	}},
	// }

	// // TODO: save event info
	// sc := bot.getServerConfig(i.GuildID)
	// _, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	// if err != nil {
	// 	log.Warn("Unable to send message %w", err)
	// }

	//
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Log message",
					Description: fmt.Sprintf("Document user behavior for for <@%v> - good, bad, or noteworthy", data.Author.ID),
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
				},
			}},
		},
	})

}

func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: bot.userInfoIntegrationresponse(i),
	})
	if err != nil {
		log.Errorf("error responding to user info (message context), err: %v", err)
	}
}
