// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
)

// An Item represents a single document in a repository.
type Item struct {
	route  *route.Route
	parent *Item
	files  []*File
}

func NewItem(path string, files []*File) (*Item, error) {

	route, err := route.New(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", path, err)
	}

	return &Item{
		route: route,
		files: files,
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
