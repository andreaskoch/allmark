// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

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
