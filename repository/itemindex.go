// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/watcher"
)

type ItemIndex struct {
	*watcher.ChangeHandler

	path  string
	items []*Item
}

func NewItemIndex(path string, items []*Item) (*ItemIndex, error) {

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(path)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for index %q.\nError: %s\n", path, err)
	}

	// create the index
	index := &ItemIndex{
		ChangeHandler: changeHandler,

		path:  path,
		items: items,
	}

	// todo update index on item change

	return index, nil
}

func (itemIndex *ItemIndex) String() string {
	return fmt.Sprintf("%s", itemIndex.path)
}

func (itemIndex *ItemIndex) Walk(walkFunc func(item *Item)) {
	for _, item := range itemIndex.items {
		item.Walk(walkFunc)
	}
}

func (itemIndex *ItemIndex) Path() string {
	return itemIndex.path
}

func (itemIndex *ItemIndex) Directory() string {
	return itemIndex.Path()
}

func (itemIndex *ItemIndex) PathType() string {
	return path.PatherTypeIndex
}
