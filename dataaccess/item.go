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
	childs []*Item
}

// Creates a new root Item that has no parent.
// A root item usually represents a repository - a collection of items.
func NewRootItem(path string, files []*File) (*Item, error) {
	return newItem(path, nil, files)
}

// Creates a new Item object that is the child of the supplied parent Item.
func NewItem(path string, parent *Item, files []*File) (*Item, error) {
	return newItem(path, parent, files)
}

func newItem(path string, parent *Item, files []*File) (*Item, error) {

	route, err := route.New(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", path, err)
	}

	return &Item{
		route:  route,
		parent: parent,
		files:  files,
	}, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) SetParent(parent *Item) {
	item.parent = parent
}

func (item *Item) Parent() *Item {
	return item.parent
}

func (item *Item) Route() *route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}

func (item *Item) SetChilds(childs []*Item) {

	// make the the current Item the parent for all childs
	for _, child := range childs {
		child.SetParent(item)
	}

	item.childs = childs
}

func (item *Item) Childs() []*Item {
	return item.childs
}

func (item *Item) Walk(callback func(item *Item)) {
	callback(item)

	for _, child := range item.childs {
		child.Walk(callback)
	}
}
