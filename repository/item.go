// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/howeyc/fsnotify"
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
	Title       string
	Description string
	RawLines    []string
	Files       *FileIndex
	MetaData    MetaData
	Type        string
	ChildItems  []*Item

	path               string
	onChangeCallbacks  map[string]func(item *Item)
	itemIsBeingWatched bool
}

func NewItem(itemPath string, childItems []*Item) (item *Item, err error) {

	// determine the type
	itemType := getItemType(itemPath)
	if itemType == UnknownItemType {
		return nil, errors.New(fmt.Sprintf("The item %q does not match any of the known item types.", itemPath))
	}

	// get the item's directory
	itemDirectory := filepath.Dir(itemPath)

	// create a new item
	item = &Item{
		ChildItems: childItems,
		Type:       itemType,
		Files:      NewFileIndex(filepath.Join(itemDirectory, FilesDirectoryName)),

		path: itemPath,
	}

	return item, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("Item %s\n", item.path)
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

	item.pauseWatch()
	walkFunc(item)
	item.resumeWatch()

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func (item *Item) RegisterOnChangeCallback(name string, callbackFunction func(item *Item)) {

	if item.onChangeCallbacks == nil {
		// initialize on first use
		item.onChangeCallbacks = make(map[string]func(item *Item))

		// start watching for changes
		item.startWatch()
	}

	if _, ok := item.onChangeCallbacks[name]; ok {
		fmt.Printf("Change callback %q already present.", name)
	}

	item.onChangeCallbacks[name] = callbackFunction
}

func (item *Item) pauseWatch() {
	fmt.Printf("Pausing watch on item %s\n", item)
	item.itemIsBeingWatched = false
}

func (item *Item) watchIsPaused() bool {
	return item.itemIsBeingWatched == false
}

func (item *Item) resumeWatch() {
	fmt.Printf("Resuming watch on item %s\n", item)
	item.itemIsBeingWatched = true
}

func (item *Item) startWatch() *Item {

	item.itemIsBeingWatched = true

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error while creating watch for item %q. Error: %v", item, err)
		return item
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:

				if !item.watchIsPaused() {
					fmt.Println("Item changed ->", event)

					for name, callback := range item.onChangeCallbacks {
						fmt.Printf("Item changed. Executing callback %q on for item %q\n", name, item)
						callback(item)
					}
				}

			case err := <-watcher.Error:
				fmt.Printf("Watch error on item %q. Error: %v\n", item, err)
			}
		}
	}()

	err = watcher.Watch(item.path)
	if err != nil {
		fmt.Printf("Error while creating watch for folder %q. Error: %v\n", item.path, err)
	}

	return item
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
