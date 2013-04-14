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
	ImageGalleryItemType = "imagegallery"
	LocationItemType     = "location"
	CommentItemType      = "comment"
	TagItemType          = "tag"
	RepositoryItemType   = "repository"

	FilesDirectoryName = "files"
)

type Item struct {
	watcher.ChangeHandler

	Title       string
	Description string
	RawLines    []string
	Files       *FileIndex
	MetaData    MetaData
	Type        string
	ChildItems  []*Item

	path string
}

func NewItem(filePath string, childItems []*Item) (item *Item, err error) {

	// determine the type
	itemType := getItemType(filePath)
	if itemType == UnknownItemType {
		return nil, fmt.Errorf("The item %q does not match any of the known item types.", filePath)
	}

	// get the item's directory
	itemDirectory := filepath.Dir(filePath)

	// create a file change handler
	fileChangeHandler, err := watcher.NewFileChangeHandler(filePath)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for item %q.\nError: %s\n", filePath, err)
	}

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, FilesDirectoryName)
	fileIndex, err := NewFileIndex(filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item = &Item{
		ChangeHandler: fileChangeHandler,
		ChildItems:    childItems,
		Type:          itemType,
		Files:         fileIndex,

		path: filePath,
	}

	reThrow := func(event *watcher.WatchEvent) {
		fmt.Println("Rethrow")
		item.Throw(event)
	}

	fileIndex.OnModify("Rethrow", reThrow)

	return item, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.path)
}

func (item *Item) Path() string {
	return item.path
}

func (item *Item) PathType() string {
	return path.PatherTypeItem
}

func (item *Item) Directory() string {
	return filepath.Dir(item.Path())
}

func (item *Item) Walk(walkFunc func(item *Item)) {

	walkFunc(item)

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func getItemType(filePath string) string {
	lowercaseFilename := strings.ToLower(filepath.Base(filePath))

	switch lowercaseFilename {
	case "document.md", "readme.md":
		return DocumentItemType

	case "presentation.md":
		return PresentationItemType

	case "collection.md":
		return CollectionItemType

	case "message.md":
		return MessageItemType

	case "imagegallery.md":
		return ImageGalleryItemType

	case "location.md":
		return LocationItemType

	case "comment.md":
		return CommentItemType

	case "tag.md":
		return TagItemType

	case "repository.md":
		return RepositoryItemType
	}

	return UnknownItemType
}
