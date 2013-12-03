// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
)

type ItemType int

const (
	TypeDocument ItemType = iota
	TypePresentation
	TypeMessage
	TypeLocation
	TypeRepository
	TypeUnknown
)

// An Item represents a single document.
type Item struct {
	route *route.Route
	files []*File

	Title       string
	Description string
	Content     string

	MetaData MetaData
}

func NewItem(route *route.Route, files []*File) (*Item, error) {
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
