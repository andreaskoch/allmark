// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"time"
)

type MetaData struct {
	Language         string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Tags             Tags
	ItemType         string
	Alias            string
	Author           string
	Locations        Locations
	GeoData          GeoInformation
}

func NewMetaData() MetaData {

	now := time.Now()

	return MetaData{
		Language:         "en",
		CreationDate:     now,
		LastModifiedDate: now,
	}
}

type GeoInformation struct {
	Street    string
	City      string
	Postcode  string
	Country   string
	Latitude  string
	Longitude string
	MapType   string
	Zoom      int
}

func NewGeoInformation() GeoInformation {
	return GeoInformation{}
}
