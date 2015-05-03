// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"allmark.io/modules/common/content"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/dataaccess"
	"fmt"

	"github.com/andreaskoch/go-fswatch"
)

// Create a new physical item.
func newPhysicalItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item, filePath string) dataaccess.Item {
	return newItem(dataaccess.TypePhysical, route, contentProvider, files, childs, filePath)
}

// Create a new virtual item.
func newVirtualItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item) dataaccess.Item {
	return newItem(dataaccess.TypeVirtual, route, contentProvider, files, childs, "")
}

// Create new file-collection item.
func newFileCollectionItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File) dataaccess.Item {
	return newItem(dataaccess.TypeFileCollection, route, contentProvider, files, nil, "")
}

// Create a new item with the given item type.
func newItem(itemType dataaccess.ItemType, route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item, filePath string) dataaccess.Item {
	return &Item{
		contentProvider,
		itemType,
		route,
		files,
		childs,

		filePath,
		nil,
	}
}

// An Item represents a single document in a repository.
type Item struct {
	*content.ContentProvider

	itemType   dataaccess.ItemType
	route      route.Route
	filesFunc  func() []dataaccess.File
	childsFunc func() []dataaccess.Item

	itemPath    string
	filewatcher *fswatch.FileWatcher
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route.String())
}

func (item *Item) Id() string {
	hash := hashutil.FromString(item.route.Value())

	return hash
}

// Get the type of this item (e.g. "physical", "virtual", ...)
func (item *Item) Type() dataaccess.ItemType {
	return item.itemType
}

// Gets a flag inidicating whether this item can have childs or not.
func (item *Item) CanHaveChilds() bool {
	switch item.Type() {

	// each child directory which is not the "files" folder can be a child
	case dataaccess.TypePhysical, dataaccess.TypeVirtual:
		return true

		// file collection items cannot have childs because all items in the directory are "files" and not items
	case dataaccess.TypeFileCollection:
		return false

	}

	panic("Unreachable. Unknown Item type.")
}

// Get the route of this item.
func (item *Item) Route() route.Route {
	return item.route
}

// Get the files of this item. Returns a slice of zero or more files.
func (item *Item) Files() (files []dataaccess.File) {

	if item.filesFunc == nil {
		return []dataaccess.File{}
	}

	return item.filesFunc()
}

func (item *Item) StartWatching() chan *Item {

	updateChannel := make(chan *Item, 1)

	item.filewatcher = fswatch.NewFileWatcher(item.itemPath, 1)
	item.filewatcher.Start()

	go func() {
		for item.filewatcher.IsRunning() {

			select {
			case <-item.filewatcher.Modified():
				updateChannel <- item
			}
		}
	}()

	return updateChannel
}

func (item *Item) StopWatching() {
	if item.filewatcher != nil {
		item.filewatcher.Stop()
	}
}
