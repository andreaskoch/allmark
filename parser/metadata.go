package parser

import (
	"andyk/docs/date"
	"andyk/docs/util"
	"strings"
	"time"
)

type MetaData struct {
	Language string
	Date     time.Time
	Tags     []string
	ItemType string
}

func (metaData MetaData) String() string {
	s := "Language: " + metaData.Language
	s += "\nDate: " + metaData.Date.String()
	s += "\nTags: " + strings.Join(metaData.Tags, ", ")
	s += "\nType: " + metaData.ItemType

	return s
}

func ParseMetaData(lines []string, itemTypeCallback func() string) (MetaData, Match, []string) {

	metaDataLocation, lines := locateMetaData(lines)
	if !metaDataLocation.Found {
		return MetaData{}, metaDataLocation, lines
	}

	// assemble meta data
	var metaData MetaData

	for _, line := range metaDataLocation.Matches {
		isKeyValuePair, matches := util.IsMatch(line, MetaDataPattern)

		// skip if line is not a key-value pair
		if !isKeyValuePair {
			continue
		}

		// prepare key and value
		key := strings.ToLower(strings.TrimSpace(matches[1]))
		value := strings.TrimSpace(matches[2])

		switch strings.ToLower(key) {

		case "language":
			{
				metaData.Language = value
				break
			}

		case "date":
			{
				date, err := date.ParseIso8601Date(value)
				if err == nil {
					metaData.Date = date
				}
				break
			}

		case "tags":
			{
				metaData.Tags = getTagsFromValue(value)
				break
			}

		case "type":
			{
				itemTypeString := strings.ToLower(value)
				if itemTypeString != "" {
					metaData.ItemType = itemTypeString
				}
				break
			}

		}
	}

	// make sure the item type is set
	if metaData.ItemType == "" {

		// use the provided item type callback
		metaData.ItemType = itemTypeCallback()
	}

	return metaData, metaDataLocation, lines
}

// locateMetaData checks if the current Document
// contains meta data.
func locateMetaData(lines []string) (Match, []string) {

	// Find the last horizontal rule in the document
	lastFoundHorizontalRulePosition := -1
	for lineNumber, line := range lines {

		lineMatchesHorizontalRulePattern := HorizontalRulePattern.MatchString(line)
		if lineMatchesHorizontalRulePattern {
			lastFoundHorizontalRulePosition = lineNumber
		}

	}

	// If there is no horizontal rule there is no meta data
	if lastFoundHorizontalRulePosition == -1 {
		return NotFound(), lines
	}

	// If the document has no more lines than
	// the last found horizontal rule there is no
	// room for meta data
	metaDataStartLine := lastFoundHorizontalRulePosition + 1
	if len(lines) <= metaDataStartLine {
		return NotFound(), lines
	}

	// Check if the last horizontal rule is followed
	// either by white space or be meta data
	for _, line := range lines[metaDataStartLine:] {

		lineMatchesMetaDataPattern := MetaDataPattern.MatchString(line)
		if lineMatchesMetaDataPattern {

			endLine := len(lines)
			return Found(lines[metaDataStartLine:endLine]), lines

		}

		lineIsEmpty := EmptyLinePattern.MatchString(line)
		if !lineIsEmpty {
			return NotFound(), lines
		}

	}

	return NotFound(), lines
}

func getTagsFromValue(value string) []string {
	rawTags := strings.Split(value, ",")
	tags := make([]string, 0, 1)

	for _, tag := range rawTags {
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			tags = append(tags, trimmedTag)
		}

	}

	return tags
}
