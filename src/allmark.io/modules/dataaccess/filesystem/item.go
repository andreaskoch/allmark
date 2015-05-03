// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/content"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/dataaccess"
	"fmt"
	"path/filepath"
)

// Create a new physical item.
func newPhysicalItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item) dataaccess.Item {
	return newItem(dataaccess.TypePhysical, route, contentProvider, files, childs)
}

// Create a new virtual item.
func newVirtualItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item) dataaccess.Item {
	return newItem(dataaccess.TypeVirtual, route, contentProvider, files, childs)
}

// Create new file-collection item.
func newFileCollectionItem(route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File) dataaccess.Item {
	return newItem(dataaccess.TypeFileCollection, route, contentProvider, files, nil)
}

// Create a new item with the given item type.
func newItem(itemType dataaccess.ItemType, route route.Route, contentProvider *content.ContentProvider, files func() []dataaccess.File, childs func() []dataaccess.Item) dataaccess.Item {
	return &Item{
		contentProvider,
		itemType,
		route,
		files,
		childs,
	}
}

// An Item represents a single document in a repository.
type Item struct {
	*content.ContentProvider

	itemType   dataaccess.ItemType
	route      route.Route
	filesFunc  func() []dataaccess.File
	childsFunc func() []dataaccess.Item
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

func newItemProvider(logger logger.Logger, repositoryPath string) (*itemProvider, error) {

	// abort if repoistory path does not exist
	if !fsutil.PathExists(repositoryPath) {
		return nil, fmt.Errorf("The repository path %q does not exist.", repositoryPath)
	}

	// abort if the supplied repository path is not a directory
	if isDirectory, _ := fsutil.IsDirectory(repositoryPath); !isDirectory {
		return nil, fmt.Errorf("The supplied item repository path %q is not a directory.", repositoryPath)
	}

	// create the file fileProvider
	provider, err := newFileProvider(logger, repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create the item provider because the file provider could not be created. Error: %s", err.Error())
	}

	return &itemProvider{
		logger:         logger,
		repositoryPath: repositoryPath,
		fileProvider:   provider,
	}, nil
}

type itemProvider struct {
	logger         logger.Logger
	repositoryPath string

	fileProvider *fileProvider
}

func (itemProvider *itemProvider) GetItemFromDirectory(itemDirectory string) (item dataaccess.Item, err error) {

	// abort if path does not exist
	if !fsutil.PathExists(itemDirectory) {
		return nil, fmt.Errorf("The path %q does not exist.", itemDirectory)
	}

	// make sure the item directory points to a folder not a file
	if isDirectory, _ := fsutil.IsDirectory(itemDirectory); !isDirectory {
		itemProvider.logger.Warn("The supplied item directory path %q is not a directory using the parent instead.", itemDirectory)
		itemDirectory = filepath.Dir(itemDirectory)
	}

	// abort if path is reserved
	if isReservedDirectory(itemDirectory) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemDirectory)
	}

	// physical item from markdown file
	if found, markdownFilePath := findMarkdownFileInDirectory(itemDirectory); found {

		// create an item from the markdown file
		return itemProvider.newItemFromFile(itemDirectory, markdownFilePath)

	}

	// virtual item
	if directoryContainsItems(itemDirectory, 3) {
		return itemProvider.newVirtualItem(itemDirectory)
	}

	// file collection item
	return itemProvider.newFileCollectionItem(itemDirectory)
}

func (itemProvider *itemProvider) getChildItemsFromDirectory(itemDirectory string) (childItems []dataaccess.Item) {

	childItems = make([]dataaccess.Item, 0)

	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		child, err := itemProvider.GetItemFromDirectory(childItemDirectory)
		if err != nil {
			itemProvider.logger.Warn("Cannot create item from directory. Error: %s", err.Error())
		}

		childItems = append(childItems, child)
	}

	return childItems
}

func (itemProvider *itemProvider) newItemFromFile(itemDirectory, filePath string) (dataaccess.Item, error) {

	route := route.NewFromItemPath(itemProvider.repositoryPath, filePath)
	itemProvider.logger.Debug("Creating a physical item from route %q", route)

	// content
	contentProvider := newFileContentProvider(filePath, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item := newPhysicalItem(route, contentProvider, files, childs)
	return item, nil
}

func (itemProvider *itemProvider) newVirtualItem(itemDirectory string) (dataaccess.Item, error) {

	route := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	itemProvider.logger.Debug("Creating a virtual item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item := newVirtualItem(route, contentProvider, files, childs)
	return item, nil
}

func (itemProvider *itemProvider) newFileCollectionItem(itemDirectory string) (dataaccess.Item, error) {

	route := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	itemProvider.logger.Debug("Creating a file collection item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := itemDirectory
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// create the item
	item := newFileCollectionItem(route, contentProvider, files)
	return item, nil
}
