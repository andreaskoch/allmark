// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common"
	"github.com/andreaskoch/allmark2/common/route"
)

// An Item represents a single document in a repository.
type Item struct {
	route *route.Route
	files []*File

	hashProvider    common.HashProviderFunc
	contentProvider common.ContentProviderFunc
}

func NewItem(route *route.Route, hashProvider common.HashProviderFunc, contentProvider common.ContentProviderFunc, files []*File) (*Item, error) {
	return &Item{
		route: route,
		files: files,

		hashProvider:    hashProvider,
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

func (item *Item) HashProvider() common.HashProviderFunc {
	return item.hashProvider
}

func (item *Item) ContentProvider() common.ContentProviderFunc {
	return item.contentProvider
}
