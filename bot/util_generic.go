package bot

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	// Enabled takes a boolean and returns "enabled" or "disabled" as a string
	Enabled = map[bool]string{
		true:  "enabled",
		false: "disabled",
	}
	// Takes a boolean and returns a colorized button if true
	ActiveButton = map[bool]discordgo.ButtonStyle{
		true:  discordgo.PrimaryButton,
		false: discordgo.SecondaryButton,
	}
)

// getFieldNamesByType takes an interface as an argument
// and returns an array of the field names. Ripped from
// https://stackoverflow.com/a/18927729
func convertFlatStructToSliceStringMap(i interface{}) []map[string]string {
	t := reflect.TypeOf(i)
	tv := reflect.ValueOf(i)

	// Keys is a list of keys of the values map
	// It's used for alphanumeric sorting later
	keys := make([]string, 0, t.NumField())

	// Values is an object that will hold an unsorted representation
	// of the interface
	values := map[string]string{}

	// Convert the struct to map[string]string
	for i := 0; i < t.NumField(); i++ {
		k := t.Field(i).Name
		v := tv.Field(i)
		values[k] = fmt.Sprintf("%v", v)
		keys = append(keys, k)
	}

	sort.Strings(keys)
	sortedValues := make([]map[string]string, 0, t.NumField())
	for _, k := range keys {
		sortedValues = append(sortedValues, map[string]string{k: values[k]})
	}

	return sortedValues
}

// getTagValue looks up the tag for a given field of the specified type
// Be advised, if the tag can't be found, it returns an empty string
func getTagValue(i interface{}, field string, tag string) string {
	r, ok := reflect.TypeOf(i).FieldByName(field)
	if !ok {
		return ""
	}
	return r.Tag.Get(tag)
}

// Returns a multiline string that pretty prints botStats
func structToPrettyDiscordFields(i any, globalMessage bool) []*discordgo.MessageEmbedField {
	var fields ([]*discordgo.MessageEmbedField)

	stringMapSlice := convertFlatStructToSliceStringMap(i)

	for _, stringMap := range stringMapSlice {
		for key, value := range stringMap {
			globalKey := getTagValue(i, key, "global") == "true"
			// If this key is a global key but
			// the message is not a global message, skip adding the field
			if globalKey && !globalMessage {
				continue
			}
			formattedKey := getTagValue(i, key, "pretty")
			newField := discordgo.MessageEmbedField{
				Name:   formattedKey,
				Value:  fmt.Sprintf("%v", value),
				Inline: getTagValue(i, key, "inline") == "",
			}
			fields = append(fields, &newField)
		}
	}

	return fields
}

// handlePlural returns the provided `src` and add `suf` if `count` is more than 1
func handlePlural(src, suf string, count int) string {
	if count > 1 {
		return src + suf
	}
	return src
}

func getUserIDFromDiscordReference(content string) string {
	pattern := regexp.MustCompile(`<@(\d+)>`)

	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	} else {
		log.Warnf("unable to get user ID from reference")
		return ""
	}
}

func getAttachmentURLs(content string) (urls []string) {
	pattern := regexp.MustCompile(`\(([^\)]+)\)`)

	matches := pattern.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				urls = append(urls, match[1])
			}
		}
	} else {
		log.Warnf("unable to get attachments ID from reference embed")
	}
	return urls
}

func nullInt64(i int) *int64 {
	v := int64(i)
	return &v
}

func nullBool(b bool) *bool {
	return &b
}
