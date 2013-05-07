// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"github.com/andreaskoch/allmark/view"
	"github.com/andreaskoch/allmark/watcher"
	"path/filepath"
)

const (
	FilesDirectoryName = "files"
)

type Item struct {
	*watcher.ChangeHandler

	ViewModel view.Model

	Files      *FileIndex
	childItems []*Item

	directory string
	path      string
	isVirtual bool
}

func NewVirtualItem(path string, childItems []*Item) (item *Item, err error) {

	if isFile, _ := util.IsFile(path); isFile {
		return nil, fmt.Errorf("Cannot create virtual items from files (%q).", path)
	}

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(path)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for item %q.\nError: %s\n", path, err)
	}

	// create the file index
	filesDirectory := filepath.Join(path, FilesDirectoryName)
	fileIndex, err := NewFileIndex(filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item = &Item{
		ChangeHandler: changeHandler,
		Files:         fileIndex,

		childItems: childItems,
		directory:  path,
		path:       path,
		isVirtual:  true,
	}

	// watch for changes in the file index
	fileIndex.OnChange("Throw Item Events on File index change", func(event *watcher.WatchEvent) {
		item.Throw(event)
	})

	return item, nil

}

func NewItem(path string, childItems []*Item) (item *Item, err error) {

	if isDirectory, _ := util.IsDirectory(path); isDirectory {
		return nil, fmt.Errorf("Cannot create items from directories (%q).", path)
	}

	// get the item's directory
	itemDirectory := filepath.Dir(path)

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(path)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for item %q.\nError: %s\n", path, err)
	}

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, FilesDirectoryName)
	fileIndex, err := NewFileIndex(filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item = &Item{
		ChangeHandler: changeHandler,
		Files:         fileIndex,

		childItems: childItems,
		directory:  itemDirectory,
		path:       path,
		isVirtual:  false,
	}

	// watch for changes in the file index
	fileIndex.OnChange("Throw Item Events on File index change", func(event *watcher.WatchEvent) {
		item.Throw(event)
	})

	return item, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.path)
}

func (item *Item) Path() string {
	return item.path
}

func (item *Item) Directory() string {
	return item.directory
}

func (item *Item) PathType() string {
	return path.PatherTypeItem
}

func (item *Item) Childs() []*Item {
	return item.childItems
}

func (item *Item) IsVirtual() bool {
	return item.isVirtual
}
