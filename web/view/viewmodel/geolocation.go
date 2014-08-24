// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type GeoLocation struct {
	PlaceName   string `json:"placename"`
	Address     string `json:"address"`
	Coordinates string `json:"coordinates"`

	Street    string `json:"street"`
	City      string `json:"city"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	MapType   string `json:"mapType"`
	Zoom      int    `json:"zoom"`
}
