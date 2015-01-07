// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"allmark.io/modules/common/content"
	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	*dataaccess.File
}

func NewFromPath(fileRoute, parentRoute route.Route, contentProvider *content.ContentProvider) (*File, error) {
	dataaccessFile, err := dataaccess.NewFile(fileRoute, parentRoute, contentProvider)
	if err != nil {
		return nil, err
	}

	return &File{dataaccessFile}, nil
}

func NewFromDataAccess(dataaccessFile *dataaccess.File) (*File, error) {
	return &File{dataaccessFile}, nil
}
