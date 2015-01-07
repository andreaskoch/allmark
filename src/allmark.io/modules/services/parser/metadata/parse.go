// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"allmark.io/modules/common/util/dateutil"
	"allmark.io/modules/model"
	"allmark.io/modules/services/parser/pattern"
	"strings"
	"time"
)

func Parse(item *model.Item, lastModifiedDate time.Time, lines []string) (parseError error) {

	// find the meta data section
	metaDataLines := GetMetaDataLines(lines)

	// create a new meta data object
	metaData := model.NewMetaData()

	// parse the different attributes
	remainingLines := parseLanguage(metaData, metaDataLines)
	remainingLines = parseAuthor(metaData, remainingLines)
	remainingLines = parseAlias(metaData, getFallbackAlias(item), remainingLines)
	remainingLines = parseCreationDate(metaData, lastModifiedDate, remainingLines)
	remainingLines = parseLastModifiedDate(metaData, lastModifiedDate, remainingLines)
	remainingLines = parseTags(metaData, remainingLines)
	remainingLines = parseGeoInformation(metaData, remainingLines)
	remainingLines = parseLocations(metaData, remainingLines)

	// assign the meta data to the item
	item.MetaData = *metaData
	return
}

func parseLanguage(metaData *model.MetaData, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"language", "lang"}, lines)
	if found {
		metaData.Language = value
	}

	return remainingLines
}

func parseAuthor(metaData *model.MetaData, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"author"}, lines)
	if found {
		metaData.Author = value
	}

	return remainingLines
}

func parseAlias(metaData *model.MetaData, fallback string, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"alias"}, lines)

	if found {
		metaData.Alias = strings.ToLower(strings.TrimSpace(value))
	} else {
		metaData.Alias = strings.ToLower(strings.TrimSpace(fallback))
	}

	return remainingLines
}

func parseTags(metaData *model.MetaData, lines []string) (remainingLines []string) {

	found, value, remainingLines := getSingleLineMetaData([]string{"tags"}, lines)

	if found {
		rawTags := strings.Split(value, ",")
		metaData.Tags, _ = model.NewTagsFromNames(rawTags)
	} else {
		// begin parser multi-line meta data
		remainingMetaDataText := strings.Join(remainingLines, "\n")

		// parse multi line tags
		if hasTags, tags := pattern.IsMultiLineTagDefinition(remainingMetaDataText); hasTags {
			metaData.Tags, _ = model.NewTagsFromNames(tags)
		}
	}

	return remainingLines
}

func parseLocations(metaData *model.MetaData, lines []string) (remainingLines []string) {

	// begin parser multi-line meta data
	remainingMetaDataText := strings.Join(lines, "\n")

	// parse multi line locations
	if hasLocations, locations := pattern.IsMultiLineLocationDefinition(remainingMetaDataText); hasLocations {
		metaData.Locations = model.NewLocationsFromNames(locations)
	}

	return lines
}

func parseCreationDate(metaData *model.MetaData, fallbackDate time.Time, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"created at", "date"}, lines)
	if found {
		date, _ := dateutil.ParseIso8601Date(value, fallbackDate)
		metaData.CreationDate = date
	} else {
		metaData.LastModifiedDate = fallbackDate
	}

	return remainingLines
}

func parseLastModifiedDate(metaData *model.MetaData, fallbackDate time.Time, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"modified at", "modified"}, lines)
	if found {
		date, _ := dateutil.ParseIso8601Date(value, fallbackDate)
		metaData.LastModifiedDate = date
	} else {
		metaData.LastModifiedDate = fallbackDate
	}

	return remainingLines
}

func getFallbackAlias(item *model.Item) string {

	route := item.Route().String()
	components := strings.Split(route, "/")

	// if the number of components is less than two than there is no item directory
	numberOfComponents := len(components)
	if numberOfComponents < 2 {
		return ""
	}

	// return the second last component as the alias
	secondLastComponentPosition := numberOfComponents - 2
	return components[secondLastComponentPosition]
}
