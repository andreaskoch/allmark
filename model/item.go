// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {
	case TypeDocument:
		return "document"

	case TypePresentation:
		return "presentation"

	case TypeMessage:
		return "message"

	case TypeLocation:
		return "location"

	case TypeRepository:
		return "repository"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

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
	route      route.Route
	files      []*File
	sourceType dataaccess.ItemType

	Type ItemType

	Title       string
	Description string
	Content     string

	MetaData *MetaData
}

func NewItem(route route.Route, files []*File, sourceType dataaccess.ItemType) *Item {

	return &Item{
		route:      route,
		files:      files,
		sourceType: sourceType,
	}

}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) FolderName() string {
	return item.route.LastComponentName()
}

func (item *Item) Route() route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}

func (item *Item) IsPhysical() bool {
	return item.sourceType == dataaccess.TypePhysical
}

func (item *Item) IsVirtual() bool {
	return item.sourceType == dataaccess.TypeVirtual
}

func (item *Item) IsFileCollection() bool {
	return item.sourceType == dataaccess.TypeFileCollection
}
