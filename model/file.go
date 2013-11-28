// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/andreaskoch/allmark2/dataaccess"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	*dataaccess.File
}

func NewFile(dataaccessFile *dataaccess.File) (*File, error) {
	return &File{dataaccessFile}, nil
}
