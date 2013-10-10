// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
	"strings"
)

func getLocations(locations repository.Locations, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) []*view.Model {
	locationModels := make([]*view.Model, 0)

	for _, location := range locations {
		item := itemResolver(location.String(), isLocation)
		if item != nil {
			locationModels = append(locationModels, getModel(item, relativePath, absolutePath, content))
		}
	}

	return locationModels
}

func getGeoLocation(item *repository.Item) *view.GeoLocation {
	return &view.GeoLocation{
		PlaceName:   getPlaceName(item),
		Address:     getAddress(item.MetaData.GeoData),
		Coordinates: getCoordinates(item.MetaData.GeoData),

		Street:    item.MetaData.GeoData.Street,
		City:      item.MetaData.GeoData.City,
		Postcode:  item.MetaData.GeoData.Postcode,
		Country:   item.MetaData.GeoData.Country,
		Latitude:  item.MetaData.GeoData.Latitude,
		Longitude: item.MetaData.GeoData.Longitude,
		MapType:   item.MetaData.GeoData.MapType,
		Zoom:      item.MetaData.GeoData.Zoom,
	}
}

func isLocation(item *repository.Item) bool {
	if item == nil {
		return false
	}

	return item.MetaData.ItemType == "location"
}

func getAddress(geoData repository.GeoInformation) string {
	components := []string{geoData.Street, geoData.Postcode, geoData.City, geoData.Country}
	return strings.Join(components, ", ")
}

func getPlaceName(item *repository.Item) string {
	if item.Title == "" || item.MetaData.GeoData.City == "" {
		return ""
	}
	components := []string{item.Title, item.MetaData.GeoData.City, item.MetaData.GeoData.Country}
	return strings.Join(components, ", ")
}

func getCoordinates(geoData repository.GeoInformation) string {
	if geoData.Latitude == "" || geoData.Longitude == "" {
		return ""
	}

	return fmt.Sprintf("%s; %s", geoData.Latitude, geoData.Longitude)
}
