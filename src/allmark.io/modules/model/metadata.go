// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// MetaData defines meta-attributes of repository items.
type MetaData struct {
	Language         string
	CreationDate     time.Time
	LastModifiedDate time.Time
	Tags             []string
	Aliases          []string
	Author           string
	GeoInformation   GeoInformation
}

// NewMetaData creates a new instance of the the MetaData struct.
func NewMetaData() *MetaData {
	return &MetaData{}
}
