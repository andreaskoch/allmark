// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/date"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/util"
	"strings"
	"time"
)

var minDate time.Time

func newMetaData(item *repository.Item) (repository.MetaData, error) {

	date, err := getItemModificationTime(item)
	if err != nil {
		return repository.MetaData{}, err
	}

	metaData := repository.MetaData{
		ItemType: types.DocumentItemType,
		Date:     date,
	}

	return metaData, nil
}

func parseMetaData(item *repository.Item, lines []string, getFallbackItemType func() string) (repository.MetaData, []string) {

	metaData := repository.MetaData{}

	// apply fallback item type
	metaData.ItemType = getFallbackItemType()

	// apply the fallback date
	fallbackDate := minDate
	if date, err := getItemModificationTime(item); err == nil {
		metaData.Date = date
	}

	// find the meta data section
	metaDataLocation, lines := locateMetaData(lines)
	if !metaDataLocation.Found {
		return metaData, lines
	}

	// parse the meta data
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

				date, _ := date.ParseIso8601Date(value, fallbackDate)
				fmt.Println(date)
				break
			}

		case "tags":
			{
				metaData.Tags = getTagsFromValue(value)
				break
			}

		case "type":
			{
				itemTypeString := strings.TrimSpace(strings.ToLower(value))
				if itemTypeString != "" {
					metaData.ItemType = itemTypeString
				}
				break
			}

		}
	}

	return metaData, lines
}

// locateMetaData checks if the current Document
// contains meta data.
func locateMetaData(lines []string) (Match, []string) {

	// Find the last horizontal rule in the document
	lastFoundHorizontalRulePosition := -1
	for lineNumber, line := range lines {
		if hrFound, _ := util.IsMatch(line, HorizontalRulePattern); hrFound {
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

	// the last line of content
	contentEndPosition := lastFoundHorizontalRulePosition - 1

	// Check if the last horizontal rule is followed
	// either by white space or be meta data
	for _, line := range lines[metaDataStartLine:] {

		lineMatchesMetaDataPattern := MetaDataPattern.MatchString(line)
		if lineMatchesMetaDataPattern {

			endLine := len(lines)
			return Found(lines[metaDataStartLine:endLine]), lines[0:contentEndPosition]

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

func getItemModificationTime(item *repository.Item) (time.Time, error) {
	date, err := util.GetModificationTime(item.Path())
	if err != nil {
		return minDate, err
	}

	return date, nil
}
