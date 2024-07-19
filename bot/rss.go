package bot

import (
	"fmt"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

const (
	// Slash commands
	ConfigureRSSFeed         = "addrss"
	RSSName                  = "rssname"
	RSSURL                   = "rssurl"
	RSSChannel               = "rsschannel"
	RSSUpdateFrequency       = "rssupdatefrequency"
	AddRSSFeed               = "addconfiguredrss"
	ListRSSFeeds             = "listrss"
	SelectRSSFeedForDeletion = "deleterss"

	// Options
	UpdateInterval = "updateinterval"
	FeedName       = "feedname"
	FeedURL        = "feedurl"
	TargetChannel  = "updatechannel"

	// Message commands
	ListRSSFeedsFromCommandContext = "List RSS Feeds"

	// Constants
	RSSFeedGoroutineInterval = time.Second * 5

	// Modals
	DeleteRSS = "selectrsstodelete"
)

var (
	ActiveRSSFeeds         = []RSSFeed{}
	UpdateFrequencyOptions = []discordgo.SelectMenuOption{
		{
			Label:   "15 minutes",
			Value:   "15",
			Default: true,
		},
		{
			Label:   "60 minutes",
			Value:   "60",
			Default: false,
		},
		{
			Label:   "4 hours",
			Value:   "240",
			Default: false,
		},
		{
			Label:   "12 hours",
			Value:   "720",
			Default: false,
		},
		{
			Label:   "24 hours",
			Value:   "1440",
			Default: false,
		},
	}
)

// Called after using the slash command, shows a modal to fill in info
// for the RSS feed
func (bot *ModeratorBot) ConfigureRSSFeedFromChatCommandContext(i *discordgo.InteractionCreate) {
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: ConfigureRSSFeed,
			Title:    "Configure RSS Feed",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    RSSName,
							Label:       "Name",
							Style:       discordgo.TextInputShort,
							Placeholder: "Name of the RSS Feed",
							Required:    true,
							MinLength:   3,
							MaxLength:   20,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    RSSURL,
							Label:       "URL",
							Style:       discordgo.TextInputShort,
							Placeholder: "URL for the RSS Feed",
							Required:    true,
							MinLength:   3,
							MaxLength:   20,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}
}

// Handles the modal input for configuring an RSS feed and starts the background
// processes to monitor it.
func (bot *ModeratorBot) ConfigureRSSFeed(i *discordgo.InteractionCreate) {
	options := i.Interaction.ApplicationCommandData().Options
	var intervalInMinutes int
	for _, o := range options {
		switch {
		default:
			log.Errorf("unknown command option %s", o.Name)
			err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: bot.generalErrorDisplayedToTheUser("Internal error"),
			})
			if err != nil {
				log.Errorf("error responding to settings interaction, err: %v", err)
			}
			return
		case o.Name == UpdateInterval:
			intervalInMinutes = int(o.IntValue())
		}
	}
	name := i.Interaction.ModalSubmitData().
		Components[0].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	url := i.Interaction.ModalSubmitData().
		Components[1].(*discordgo.ActionsRow).
		Components[0].(*discordgo.TextInput).
		Value

	var channelId string
	selectedChannels := i.Interaction.ApplicationCommandData().Resolved.Channels
	if len(selectedChannels) > 1 {
		log.Errorf("multiple channels selected")
		return
	}
	for _, selectedChannel := range selectedChannels {
		channelId = selectedChannel.ID
	}

	err := bot.AddRSSFeed(RSSFeed{
		Name:            name,
		URL:             url,
		UpdateFrequency: time.Duration(time.Minute * time.Duration(intervalInMinutes)),
		GuildConfigID:   i.GuildID,
		TargetChannelID: channelId,
	})
	if err != nil {
		err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: bot.generalErrorDisplayedToTheUser(fmt.Sprintf("Error registering RSS Feed: %s", err)),
		})
		if err != nil {
			log.Errorf("error responding to settings interaction, err: %v", err)
		}
		return
	}
	err = bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Title:   "Successfully registered RSS Feed",
			Content: fmt.Sprintf("**Name**: %s\n**URL**: %s", name, url),
		},
	})

	if err != nil {
		log.Errorf("error showing custom slash command creation modal, err: %v", err)
	}
}

// Creates an RSS feed goroutine that manages the feed
func (bot *ModeratorBot) AddRSSFeed(newfeed RSSFeed) (err error) {
	existingMatchingFeeds := int64(0)
	bot.DB.Where(RSSFeed{URL: newfeed.URL,
		Name:          newfeed.Name,
		GuildConfigID: newfeed.GuildConfigID}).Count(&existingMatchingFeeds)
	if existingMatchingFeeds > 0 {
		return fmt.Errorf("rss feed already exists")
	}
	// Register the feed in the database
	var registeredFeed RSSFeed
	bot.DB.Model(RSSFeed{}).Updates(newfeed).First(&registeredFeed)
	go bot.WatchFeed(registeredFeed)
	return nil
}

// Starts the RSS feed watcher which checks for updates to the RSS feed
// according to user-defined configuration
func (bot *ModeratorBot) WatchFeed(feed RSSFeed) {
	checkAfter := time.Now().Add(feed.UpdateFrequency)
	for {
		// Check if the feed is still in the list of active RSS Feeds
		idx := slices.IndexFunc(ActiveRSSFeeds, func(r RSSFeed) bool {
			return r == feed
		})
		active := idx != -1
		if !active {
			log.Infof("feed is no longer active: %s - %s", feed.Name, feed.URL)
			return
		}

		if time.Now().Unix() >= checkAfter.Unix() {
			// We should now post an update for the RSS feed
			fp := gofeed.NewParser()
			remoteFeed, err := fp.ParseURL(feed.URL)
			if err != nil {
				message := fmt.Sprintf("error checking feed %s(%s), %s", feed.Name, feed.URL, err)
				log.Error(message)
				bot.generalErrorDisplayedToTheUser(message + "\nIf this seems " +
					"like and issue with this bot, please report it!")
			}

			var author discordgo.MessageEmbedAuthor
			var thumbnail discordgo.MessageEmbedThumbnail
			var video discordgo.MessageEmbedVideo
			var image discordgo.MessageEmbedImage
			if feed.ShowAuthor {
				var authornames string
				for i, author := range remoteFeed.Authors {
					authornames = authornames + fmt.Sprintf("%s", author.Name)
					if i < len(remoteFeed.Authors)-1 {
						authornames = authornames + ", "
					}
				}
				author = discordgo.MessageEmbedAuthor{ // TODO
					Name: authornames,
				}
			}
			if feed.ShowImage {
				if feed.ShowThumbnail {
					thumbnail = discordgo.MessageEmbedThumbnail{ // TODO
						URL: remoteFeed.Image.URL,
					}
				} else {
					image = discordgo.MessageEmbedImage{
						URL: remoteFeed.Image.URL,
					}
				}
			}
			if feed.ShowVideo {
				video = discordgo.MessageEmbedVideo{ // TODO
					URL: remoteFeed.Image.URL,
				}
			}
			ms := discordgo.MessageSend{ // TODO
				Embeds: []*discordgo.MessageEmbed{{
					Title:       remoteFeed.Title,
					Description: remoteFeed.Description,
					Timestamp:   remoteFeed.Published,
					Author:      &author,
					Thumbnail:   &thumbnail,
					Video:       &video,
					Image:       &image,
					Provider: &discordgo.MessageEmbedProvider{
						Name: feed.Name,
						URL:  feed.URL,
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Update this feed or others with /%s", ListRSSFeeds),
					},
				}},
			}
			_, err = bot.DG.ChannelMessageSendComplex(feed.TargetChannelID, &ms)
			if err != nil {
				log.Error("unable to send RSS item message: %w", err)
				return
			}
			checkAfter = time.Now().Add(feed.UpdateFrequency)
		}

		// This is much less than the update frequency so we don't have to wait
		// a whole update frequency interval (at least 15 minutes) in order
		// to apply updates to the feed config
		time.Sleep(RSSFeedGoroutineInterval)
	}
}

// Starts the RSS Feed operator which manages adding and removing feeds (and
// their respective feed watching goroutines)
func (bot *ModeratorBot) StartWatchingRegisteredFeeds(guildId string) {
	// Manage feeds in DB first
	var feeds []RSSFeed
	bot.DB.Model(RSSFeed{}).Where(RSSFeed{GuildConfigID: guildId}).Find(&feeds)
	for _, feed := range feeds {
		go bot.WatchFeed(feed)
	}
}

// Lists all configured RSS Feeds for that given guild
func (bot *ModeratorBot) ListRSSFeedsFromChatCommandContext(i *discordgo.InteractionCreate) {
	// TODO: allow modifying feeds from this list
	var RSSFeeds []RSSFeed
	_ = bot.DB.Where(RSSFeed{GuildConfigID: i.GuildID}).Find(&RSSFeeds)
	var content string
	if len(RSSFeeds) != 0 {
		for _, feed := range RSSFeeds {
			updateFrequency := getHumanFriendlyLabelForInterval(feed.UpdateFrequency)
			content += fmt.Sprintf("**%s**: %s (%s)\n", feed.Name, feed.URL, updateFrequency)
		}
	} else {
		content = "There are no configured RSS feeds. Configure one with `/" + ConfigureRSSFeed + "`."
	}
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: ListRSSFeedsFromCommandContext,
			Flags:    discordgo.MessageFlagsEphemeral,
			Content:  content,
		},
	})
	if err != nil {
		log.Errorf("error responding to RSS list interaction, err: %v", err)
	}
}

// Deletes one RSS Feed
func (bot *ModeratorBot) DeleteRSSFeedFromChatCommandContext(i *discordgo.InteractionCreate) {
	var feeds []RSSFeed
	var options []discordgo.SelectMenuOption
	bot.DB.Model(RSSFeed{}).Where(RSSFeed{GuildConfigID: i.GuildID}).Find(&feeds)
	for _, feed := range feeds {
		options = append(options, discordgo.SelectMenuOption{
			Label:       feed.Name,
			Value:       feed.ID,
			Description: feed.URL,
		})
	}
	err := bot.DG.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: DeleteRSS,
			Title:    "Select RSS Feed to delete",
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    SelectRSSFeedForDeletion,
						Placeholder: "RSS Feed to delete",
						Options:     options,
					},
				},
			}},
		},
	})
	if err != nil {
		log.Warn("error responding to interaction: %w", err)
	}

}

// Converts an interval in minutes to a human friendly representation by
// using the options presented when choosing the interval
func getHumanFriendlyLabelForInterval(duration time.Duration) string {
	minutes := time.Duration.Minutes(duration)
	idx := slices.IndexFunc(UpdateFrequencyOptions, func(o discordgo.SelectMenuOption) bool { return o.Value == fmt.Sprintf("%v", minutes) })
	if idx != -1 {
		return UpdateFrequencyOptions[idx].Value
	} else {
		log.Warnf("unknown minute interval %v", minutes)
		return "Unknown"
	}
}
