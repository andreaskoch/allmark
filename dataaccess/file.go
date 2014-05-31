// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	*content.ContentProvider

	parentRoute *route.Route
	fileRoute   *route.Route
}

func NewFile(fileRoute, parentRoute *route.Route, contentProvider *content.ContentProvider) (*File, error) {
	return &File{
		contentProvider,
		parentRoute,
		fileRoute,
	}, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.fileRoute.Value())
}

func (file *File) Parent() *route.Route {
	return file.parentRoute
}

func (file *File) Route() *route.Route {
	return file.fileRoute
}
