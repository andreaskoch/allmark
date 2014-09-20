// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

type MetaData struct {
	Language         string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Tags             Tags
	Alias            string
	Author           string
	Locations        Locations
	GeoInformation   GeoInformation
}

func NewMetaData() *MetaData {

	return &MetaData{
		Language: "en",
	}
}
