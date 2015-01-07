// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"strings"
)

func getGeoLocation(item *model.Item) viewmodel.GeoLocation {
	emptyLocation := model.GeoInformation{}
	if item.MetaData.GeoInformation == emptyLocation {
		return viewmodel.GeoLocation{}
	}

	return viewmodel.GeoLocation{
		PlaceName:   getPlaceName(item),
		Address:     getAddress(item.MetaData.GeoInformation),
		Coordinates: getCoordinates(item.MetaData.GeoInformation),

		Street:    item.MetaData.GeoInformation.Street,
		City:      item.MetaData.GeoInformation.City,
		Postcode:  item.MetaData.GeoInformation.Postcode,
		Country:   item.MetaData.GeoInformation.Country,
		Latitude:  item.MetaData.GeoInformation.Latitude,
		Longitude: item.MetaData.GeoInformation.Longitude,
		MapType:   item.MetaData.GeoInformation.MapType,
		Zoom:      item.MetaData.GeoInformation.Zoom,
	}
}

func getAddress(geoData model.GeoInformation) string {
	components := []string{}

	if geoData.Street != "" {
		components = append(components, geoData.Street)
	}

	if geoData.Postcode != "" {
		components = append(components, geoData.Postcode)
	}

	if geoData.City != "" {
		components = append(components, geoData.City)
	}

	if geoData.Country != "" {
		components = append(components, geoData.Country)
	}

	if len(components) > 0 {
		return strings.Join(components, ", ")
	}

	return ""
}

func getPlaceName(item *model.Item) string {
	if item.Title == "" || item.MetaData.GeoInformation.City == "" {
		return ""
	}

	components := []string{}

	if item.Title != "" {
		components = append(components, item.Title)
	}

	if item.MetaData.GeoInformation.City != "" {
		components = append(components, item.MetaData.GeoInformation.City)
	}

	if item.MetaData.GeoInformation.Country != "" {
		components = append(components, item.MetaData.GeoInformation.Country)
	}

	if len(components) > 0 {
		return strings.Join(components, ", ")
	}

	return ""
}

func getCoordinates(geoData model.GeoInformation) string {
	if geoData.Latitude == "" || geoData.Longitude == "" {
		return ""
	}

	return fmt.Sprintf("%s; %s", geoData.Latitude, geoData.Longitude)
}
