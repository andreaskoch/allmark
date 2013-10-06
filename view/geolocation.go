// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

type GeoLocation struct {
	Street    string `json:"street"`
	City      string `json:"city"`
	Postcode  string `json:"postcode"`
	Country   string `json:"country"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	MapType   string `json:"mapType"`
	Zoom      int    `json:"zoom"`
}
