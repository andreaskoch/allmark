// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"github.com/andreaskoch/allmark2/model"
	"strconv"
)

func parseGeoInformation(metaData *model.MetaData, lines []string) (remainingLines []string) {

	geoData := &metaData.GeoInformation

	remainingLines = parseStreet(geoData, lines)
	remainingLines = parseCity(geoData, lines)
	remainingLines = parseCountry(geoData, lines)
	remainingLines = parseLatituide(geoData, lines)
	remainingLines = parseLongitude(geoData, lines)
	remainingLines = parseMapType(geoData, lines)
	remainingLines = parseZoom(geoData, lines)

	metaData.GeoInformation = *geoData

	return remainingLines
}

func parseStreet(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"street"}, lines)
	if found {
		geoInformation.Street = value
	}

	return remainingLines
}

func parseCity(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"city"}, lines)
	if found {
		geoInformation.City = value
	}

	return remainingLines
}

func parseCountry(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"country"}, lines)
	if found {
		geoInformation.Country = value
	}

	return remainingLines
}

func parseLatituide(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"latitude", "lat"}, lines)
	if found {
		geoInformation.Latitude = value
	}

	return remainingLines
}

func parseLongitude(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"longitude", "long"}, lines)
	if found {
		geoInformation.Longitude = value
	}

	return remainingLines
}

func parseMapType(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"maptype"}, lines)
	if found {
		geoInformation.Longitude = value
	}

	return remainingLines
}

func parseZoom(geoInformation *model.GeoInformation, lines []string) (remainingLines []string) {
	found, value, remainingLines := getSingleLineMetaData([]string{"zoom"}, lines)
	if found {
		if zoomLevel, err := strconv.ParseInt(value, 10, 0); err != nil && zoomLevel >= 0 && zoomLevel <= 100 {
			geoInformation.Zoom = int(zoomLevel)
		} else {
			geoInformation.Zoom = 75
		}
	}

	return remainingLines
}
