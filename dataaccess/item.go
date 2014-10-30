// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {

	case TypePhysical:
		return "physical"

	case TypeVirtual:
		return "virtual"

	case TypeFileCollection:
		return "filecollection"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

const (
	TypePhysical ItemType = iota
	TypeVirtual
	TypeFileCollection
)

// An Item represents a single document in a repository.
type Item struct {
	*content.ContentProvider
	itemType   ItemType
	route      route.Route
	filesFunc  func() []*File
	childsFunc func() []*Item
}

func NewPhysicalItem(route route.Route, contentProvider *content.ContentProvider, files func() []*File, childs func() []*Item) (*Item, error) {
	return newItem(TypePhysical, route, contentProvider, files, childs)
}

func NewVirtualItem(route route.Route, contentProvider *content.ContentProvider, files func() []*File, childs func() []*Item) (*Item, error) {
	return newItem(TypeVirtual, route, contentProvider, files, childs)
}

func NewFileCollectionItem(route route.Route, contentProvider *content.ContentProvider, files func() []*File) (*Item, error) {
	return newItem(TypeFileCollection, route, contentProvider, files, nil)
}

func newItem(itemType ItemType, route route.Route, contentProvider *content.ContentProvider, files func() []*File, childs func() []*Item) (*Item, error) {
	return &Item{
		contentProvider,
		itemType,
		route,
		files,
		childs,
	}, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route.String())
}

// Get the type of this item (e.g. "physical", "virtual", ...)
func (item *Item) Type() ItemType {
	return item.itemType
}

// Gets a flag inidicating whether this item can have childs or not.
func (item *Item) CanHaveChilds() bool {
	switch item.Type() {

	// each child directory which is not the "files" folder can be a child
	case TypePhysical, TypeVirtual:
		return true

		// file collection items cannot have childs because all items in the directory are "files" and not items
	case TypeFileCollection:
		return false

	}

	panic("Unreachable. Unknown Item type.")
}

// Get the route of this item.
func (item *Item) Route() route.Route {
	return item.route
}

// Get the childs of this item. Returns nil if this item cannot have childs; otherwise returns a slice with zero or more childs.
func (item *Item) GetChilds() (childs []*Item) {
	if !item.CanHaveChilds() || item.childsFunc == nil {
		return
	}

	return item.childsFunc()
}

// Get the files of this item. Returns a slice of zero or more files.
func (item *Item) Files() (files []*File) {

	if item.filesFunc == nil {
		return []*File{}
	}

	return item.filesFunc()
}
