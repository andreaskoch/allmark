// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
)

type ItemIndex struct {
	path  string
	items []*Item
}

func NewItemIndex(directory string, items []*Item) *ItemIndex {
	return &ItemIndex{
		path:  directory,
		items: items,
	}
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
