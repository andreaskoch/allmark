// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	*dataaccess.File
}

func NewFromPath(fileRoute, parentRoute *route.Route, contentProvider *content.ContentProvider) (*File, error) {
	dataaccessFile, err := dataaccess.NewFile(fileRoute, parentRoute, contentProvider)
	if err != nil {
		return nil, err
	}

	return &File{dataaccessFile}, nil
}

func NewFromDataAccess(dataaccessFile *dataaccess.File) (*File, error) {
	return &File{dataaccessFile}, nil
}
