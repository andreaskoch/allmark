// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"github.com/andreaskoch/allmark2/common/util/dateutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/pattern"
	"strconv"
	"strings"
	"time"
)

func Parse(item *model.Item, lines []string) (parseError error) {

	metaData := model.NewMetaData()

	// fallback: alias
	metaData.Alias = getFallbackAlias(item)

	// find the meta data section
	metaDataLines := GetLines(lines)

	// parse the single line meta data
	remainingLines := make([]string, 0)
	for _, line := range metaDataLines {

		key, value := pattern.GetSingleLineMetaDataKeyAndValue(line)

		// skip if line is not a key-value pair
		if key == "" && value == "" {
			remainingLines = append(remainingLines, line)
			continue
		}

		// prepare key and value
		key = strings.ToLower(key)
		value = strings.TrimSpace(value)

		switch key {

		case "language":
			{
				metaData.Language = value
				break
			}

		case "created at", "date":
			{
				date, _ := dateutil.ParseIso8601Date(value, time.Now())
				metaData.CreationDate = date
				break
			}

		case "modified at":
			{
				date, _ := dateutil.ParseIso8601Date(value, time.Now())
				metaData.LastModifiedDate = date
				break
			}

		case "tags":
			{
				if strings.TrimSpace(value) != "" {
					metaData.Tags = getTagsFromValue(value)
				}
				break
			}

		case "type":
			{
				if typeValue := strings.TrimSpace(strings.ToLower(value)); typeValue != "" {
					metaData.ItemType = typeValue
				}
				break
			}

		case "alias":
			{
				if aliasValue := strings.TrimSpace(strings.ToLower(value)); aliasValue != "" {
					metaData.Alias = aliasValue
				}
				break
			}

		case "author":
			{
				if author := strings.TrimSpace(value); author != "" {
					metaData.Author = author
				}
				break
			}

		case "street":
			{
				if street := strings.TrimSpace(value); street != "" {
					metaData.GeoData.Street = street
				}
				break
			}

		case "city":
			{
				if city := strings.TrimSpace(value); city != "" {
					metaData.GeoData.City = city
				}
				break
			}

		case "postcode":
			{
				if postcode := strings.TrimSpace(value); postcode != "" {
					metaData.GeoData.Postcode = postcode
				}
				break
			}

		case "country":
			{
				if country := strings.TrimSpace(value); country != "" {
					metaData.GeoData.Country = country
				}
				break
			}

		case "latitude":
			{
				if latitude := strings.TrimSpace(value); latitude != "" {
					metaData.GeoData.Latitude = latitude
				}
				break
			}

		case "longitude":
			{
				if longitude := strings.TrimSpace(value); longitude != "" {
					metaData.GeoData.Longitude = longitude
				}
				break
			}

		case "maptype":
			{
				if maptype := strings.TrimSpace(value); maptype != "" {
					metaData.GeoData.MapType = maptype
				}
				break
			}

		case "zoom":
			{
				if zoom := strings.TrimSpace(value); zoom != "" {
					if zoomLevel, err := strconv.ParseInt(zoom, 10, 0); err != nil && zoomLevel >= 0 && zoomLevel <= 100 {
						metaData.GeoData.Zoom = int(zoomLevel)
					} else {
						metaData.GeoData.Zoom = 75
					}
				}
				break
			}

		}
	}

	// begin parsing multi-line meta data
	remainingMetaDataText := strings.Join(remainingLines, "\n")

	// parse multi line tags
	if hasTags, tags := pattern.IsMultiLineTagDefinition(remainingMetaDataText); hasTags {
		metaData.Tags = model.NewTagsFromNames(tags)
	}

	// parse multi line locations
	if hasLocations, locations := pattern.IsMultiLineLocationDefinition(remainingMetaDataText); hasLocations {
		metaData.Locations = model.NewLocationsFromNames(locations)
	}

	item.MetaData = metaData
	return

}

func getFallbackAlias(item *model.Item) string {
	return "fallback alias"
}

func getTagsFromValue(value string) model.Tags {
	rawTags := strings.Split(value, ",")
	return model.NewTagsFromNames(rawTags)
}
