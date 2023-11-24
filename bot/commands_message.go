package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/tyzbit/go-discord-modtools/globals"
)

// Moderation modal menu
func (bot *ModeratorBot) ModerateFromMessageContext(i *discordgo.InteractionCreate) {
	message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	if message.ID == "" {
		log.Warn("message was provided")
	}

	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: globals.RespondToModerationModalFromMessageContext,
			Title:    "Moderate " + i.Member.User.Username,
			Embeds: []*discordgo.MessageEmbed{{
				Title: "Embed!",
				Fields: []*discordgo.MessageEmbedField{{
					Name:  "Field1!",
					Value: "Value!",
				}},
			}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder: "User",
							MenuType:    discordgo.UserSelectMenu,
							CustomID:    i.Member.User.ID,
							Options: []discordgo.SelectMenuOption{{
								Label: "UserID",
								Value: "Saved",
							}},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder:  "Channel",
							MenuType:     discordgo.ChannelSelectMenu,
							ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
							CustomID:     message.ChannelID,
							Options: []discordgo.SelectMenuOption{{
								Label: "Channel",
								Value: "Saved",
							}},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							Placeholder: "Message",
							CustomID:    message.ID,
							Options: []discordgo.SelectMenuOption{{
								Label: "Message",
								Value: "Saved",
							}},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    globals.ReasonOption,
							Label:       "Reason",
							Style:       discordgo.TextInputShort,
							Placeholder: "Why are you moderating this user or content?",
							Required:    true,
							MinLength:   1,
							MaxLength:   500,
						},
					},
				},
			},
		},
	})
}

func (bot *ModeratorBot) GetUserInfoFromMessageContext(i *discordgo.InteractionCreate) {
	return
}

func (bot *ModeratorBot) SaveEvidenceFromMessageContext(i *discordgo.InteractionCreate) {
	message := *i.Interaction.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]

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
	ms := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title: "Moder8s Evidence",
			Description: fmt.Sprintf("Collected by <@%s> on %s (<t:%v:R>)",
				i.Interaction.Member.User.ID,
				time.Now().Format(time.RFC1123Z),
				time.Now().Unix(),
			),
			Fields: fields,
			Color:  globals.Purple,
		}},
	}

	sc := bot.getServerConfig(i.GuildID)
	_, err := bot.DG.ChannelMessageSendComplex(sc.EvidenceChannelSettingID, &ms)
	if err != nil {
		log.Warn("Unable to send message %w", err)
	}

	// TODO: more information
	_ = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Evidence saved to <#" + sc.EvidenceChannelSettingID + ">",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
