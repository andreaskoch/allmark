// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
)

// An Item represents a single document in a repository.
type Item struct {
	path   string
	parent *Item
	files  []*File
	childs []*Item
}

// Creates a new root Item that has no parent.
// A root item usually represents a repository - a collection of items.
func NewRootItem(path string, files []*File, childs []*Item) *Item {
	return newItem(path, nil, files, childs)
}

// Creates a new Item object that is the child of the supplied parent Item.
func NewItem(path string, parent *Item, files []*File, childs []*Item) *Item {
	return newItem(path, parent, files, childs)
}

func newItem(path string, parent *Item, files []*File, childs []*Item) *Item {

	normalizedPath := NormalizePath(path)
	if normalizedPath == "" {
		panic("An item path cannot be empty.")
	}

	return &Item{
		path:   normalizedPath,
		parent: parent,
		files:  files,
		childs: childs,
	}
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.path)
}

func (item *Item) Parent() *Item {
	return item.parent
}

func (item *Item) Path() string {
	return item.path
}

func (item *Item) Files() []*File {
	return item.files
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
