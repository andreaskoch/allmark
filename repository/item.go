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

	*view.Model

	Files  *FileIndex
	Childs []*Item

	directory string
	path      string
	isVirtual bool

	pathProvider *path.Provider
}

func NewVirtualItem(itemPath string, childItems []*Item, pathProvider *path.Provider) (item *Item, err error) {

	if isFile, _ := util.IsFile(itemPath); isFile {
		return nil, fmt.Errorf("Cannot create virtual items from files (%q).", itemPath)
	}

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(itemPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for item %q.\nError: %s\n", itemPath, err)
	}

	// create the file index
	filesDirectory := filepath.Join(itemPath, FilesDirectoryName)
	fileIndex, err := NewFileIndex(filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item = &Item{
		ChangeHandler: changeHandler,
		Childs:        childItems,
		Files:         fileIndex,

		pathProvider: pathProvider,
		directory:    itemPath,
		path:         itemPath,
		isVirtual:    true,
	}

	// watch for changes in the file index
	fileIndex.OnChange("Throw Item Events on File index change", func(event *watcher.WatchEvent) {
		item.Throw(event)
	})

	return item, nil

}

func NewItem(itemPath string, childItems []*Item, pathProvider *path.Provider) (item *Item, err error) {

	if isDirectory, _ := util.IsDirectory(itemPath); isDirectory {
		return nil, fmt.Errorf("Cannot create items from directories (%q).", itemPath)
	}

	// get the item's directory
	itemDirectory := filepath.Dir(itemPath)

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(itemPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for item %q.\nError: %s\n", itemPath, err)
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
		Childs:        childItems,
		Files:         fileIndex,

		pathProvider: pathProvider,
		directory:    itemDirectory,
		path:         itemPath,
		isVirtual:    false,
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

func (item *Item) IsVirtual() bool {
	return item.isVirtual
}

func (item *Item) PathProvider() *path.Provider {
	return item.pathProvider
}
