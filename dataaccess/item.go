// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
)

// An Item represents a single document in a repository.
type Item struct {
	route           *route.Route
	files           []*File
	contentProvider *content.ContentProvider
}

func NewItem(route *route.Route, contentProvider *content.ContentProvider, files []*File) (*Item, error) {
	return &Item{
		route:           route,
		files:           files,
		contentProvider: contentProvider,
	}, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) Route() *route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}

func (item *Item) ContentProvider() *content.ContentProvider {
	return item.contentProvider
}
