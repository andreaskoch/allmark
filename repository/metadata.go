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
}
