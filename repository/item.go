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
	*watcher.FileWatcher

	Title       string
	Description string
	RawLines    []string
	Files       *FileIndex
	MetaData    MetaData
	Type        string
	ChildItems  []*Item

	path              string
	onChangeCallbacks map[string]func(event *watcher.WatchEvent)
}

func NewItem(itemPath string, childItems []*Item) (item *Item, err error) {

	// determine the type
	itemType := getItemType(itemPath)
	if itemType == UnknownItemType {
		return nil, fmt.Errorf("The item %q does not match any of the known item types.", itemPath)
	}

	// get the item's directory
	itemDirectory := filepath.Dir(itemPath)

	// create a watcher
	watcher, err := watcher.NewFileWatcher(itemPath)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to create a watch for item %q. Error: %s", item, err)
	}

	// create the file index
	fileIndex := NewFileIndex(filepath.Join(itemDirectory, FilesDirectoryName))

	// create the item
	item = &Item{
		FileWatcher: watcher,
		ChildItems:  childItems,
		Type:        itemType,
		Files:       fileIndex,

		path: itemPath,
	}

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

	item.Pause()
	walkFunc(item)
	item.Resume()

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func (item *Item) RegisterOnChangeCallback(name string, callbackFunction func(event *watcher.WatchEvent)) {

	if item.onChangeCallbacks == nil {
		// initialize on first use
		item.onChangeCallbacks = make(map[string]func(event *watcher.WatchEvent))

		// start watching for changes
		go func() {
			for {
				select {
				case event := <-item.Event:

					fmt.Printf("%s: %s\n", strings.ToUpper(event.Type.String()), event.Filepath)
					for _, callback := range item.onChangeCallbacks {

						item.Pause()
						callback(event)
						item.Resume()

					}
				}
			}
		}()
	}

	if _, ok := item.onChangeCallbacks[name]; ok {
		fmt.Printf("Change callback %q already present.", name)
	}

	item.onChangeCallbacks[name] = callbackFunction
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
