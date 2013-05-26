// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
)

type ItemIndex struct {
	path string
	root *Item
}

func NewItemIndex(indexPath string) (*ItemIndex, error) {

	// check if path exists
	if !util.PathExists(indexPath) {
		return nil, fmt.Errorf("The path %q does not exist.", indexPath)
	}

	if isReservedDirectory(indexPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an index.", indexPath)
	}

	// check if the path is a directory
	if isDirectory, _ := util.IsDirectory(indexPath); !isDirectory {
		indexPath = filepath.Dir(indexPath)
	}

	rootItem, err := NewItem(indexPath, 0)
	if err != nil {
		return nil, err
	}

	// create the index
	index := &ItemIndex{
		path: indexPath,
		root: rootItem,
	}

	return index, nil
}

func (itemIndex *ItemIndex) String() string {
	return fmt.Sprintf("%s", itemIndex.path)
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

func (itemIndex *ItemIndex) Root() *Item {
	return itemIndex.root
}
