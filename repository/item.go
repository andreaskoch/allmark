// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/watcher"
	"path/filepath"
	"strings"
)

const (
	UnknownItemType      = "unknown"
	DocumentItemType     = "document"
	PresentationItemType = "presentation"
	CollectionItemType   = "collection"
	MessageItemType      = "message"
	RepositoryItemType   = "repository"

	FilesDirectoryName = "files"
)

type Item struct {
	*watcher.ChangeHandler

	Type       string
	Files      *FileIndex
	childItems []*Item

	path string
}

func NewItem(path string, childItems []*Item) (item *Item, err error) {

	// determine the type
	itemType := getItemType(path)

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
		Type:          itemType,
		Files:         fileIndex,

		childItems: childItems,
		path:       path,
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
	return filepath.Dir(item.Path())
}

func (item *Item) PathType() string {
	return path.PatherTypeItem
}

func (item *Item) Childs() []*Item {
	return item.childItems
}

func getItemType(filePath string) string {
	extension := filepath.Ext(filePath)
	filenameWithExtension := filepath.Base(filePath)
	filename := filenameWithExtension[0:(strings.LastIndex(filenameWithExtension, extension))]

	switch strings.ToLower(filename) {
	case DocumentItemType:
		return DocumentItemType

	case PresentationItemType:
		return PresentationItemType

	case CollectionItemType:
		return CollectionItemType

	case MessageItemType:
		return MessageItemType

	case RepositoryItemType:
		return RepositoryItemType
	}

	return UnknownItemType
}
