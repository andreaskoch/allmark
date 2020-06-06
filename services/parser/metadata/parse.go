// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"regexp"
	"strings"
	"time"

	"github.com/elWyatt/allmark/common/util/dateutil"
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/services/parser/pattern"
)

var aliasForbiddenCharacters = regexp.MustCompile(`[^\w\d-_]`)

// Parse parses the supplied lines and writes the result to the specified item.
func Parse(item *model.Item, lastModifiedDate time.Time, lines []string) (parseError error) {

	// find the meta data section
	metaDataLines := GetMetaDataLines(lines)

	// create a new meta data object
	metaData := model.NewMetaData()

	// parse the different attributes
	remainingLines := parseLanguage(metaData, metaDataLines)
	remainingLines = parseAuthor(metaData, remainingLines)
	remainingLines = parseAlias(metaData, remainingLines)
	remainingLines = parseCreationDate(metaData, lastModifiedDate, remainingLines)
	remainingLines = parseLastModifiedDate(metaData, lastModifiedDate, remainingLines)
	remainingLines = parseTags(metaData, remainingLines)
	remainingLines = parseGeoInformation(metaData, remainingLines)

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

func parseAlias(metaData *model.MetaData, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"alias"}, lines)

	if found {

		rawAliases := strings.Split(value, ",")
		metaData.Aliases = normalizeAliases(rawAliases)

	} else {

		// begin parser multi-line meta data
		remainingMetaDataText := strings.Join(remainingLines, "\n")

		// parse multi line tags
		if hasAliases, rawAliases := pattern.IsMultiLineAliasDefinition(remainingMetaDataText); hasAliases {
			metaData.Aliases = normalizeAliases(rawAliases)
		}

	}

	return remainingLines
}

func parseTags(metaData *model.MetaData, lines []string) (remainingLines []string) {

	found, value, remainingLines := getSingleLineMetaData([]string{"tags"}, lines)

	if found {
		rawTags := strings.Split(value, ",")
		metaData.Tags = normalizeTags(rawTags)
	} else {
		// begin parser multi-line meta data
		remainingMetaDataText := strings.Join(remainingLines, "\n")

		// parse multi line tags
		if hasTags, tags := pattern.IsMultiLineTagDefinition(remainingMetaDataText); hasTags {
			metaData.Tags = normalizeTags(tags)
		}
	}

	return remainingLines
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

// normalizeAliases normalizes the given list of raw aliases.
func normalizeAliases(rawAliases []string) []string {
	var normalizedAliases []string
	for _, rawAlias := range rawAliases {

		normalizedAlias := normalizeAlias(rawAlias)

		// skip empty values
		if normalizedAlias == "" {
			continue
		}

		normalizedAliases = append(normalizedAliases, normalizedAlias)
	}

	return normalizedAliases
}

// normalizeAlias normalizes any given alias and replaces invalid characters.
func normalizeAlias(rawAlias string) string {

	// Trim whitespace
	cleaned := strings.TrimSpace(rawAlias)

	// lowercase
	cleaned = strings.ToLower(cleaned)

	// Replace spaces with dashes
	cleaned = strings.Replace(cleaned, " ", "-", -1)

	// Replace forbidden characters
	cleaned = aliasForbiddenCharacters.ReplaceAllString(cleaned, "")

	return cleaned
}

// normalizeTags normalizes the given list of raw tags.
func normalizeTags(rawTags []string) []string {
	var normalizedTags []string
	for _, rawTag := range rawTags {

		normalizedTag := normalizeTagName(rawTag)

		// skip empty values
		if normalizedTag == "" {
			continue
		}

		normalizedTags = append(normalizedTags, normalizedTag)
	}

	return normalizedTags
}

// normalizeTagName returns a normalized version of the given raw tag name.
func normalizeTagName(rawTagName string) string {
	return strings.TrimSpace(rawTagName)
}
