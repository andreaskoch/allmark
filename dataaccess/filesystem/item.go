// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"path/filepath"
)

func newItemProvider(logger logger.Logger, repositoryPath string) (*itemProvider, error) {

	// abort if repoistory path does not exist
	if !fsutil.PathExists(repositoryPath) {
		return nil, fmt.Errorf("The repository path %q does not exist.", repositoryPath)
	}

	// abort if the supplied repository path is not a directory
	if isDirectory, _ := fsutil.IsDirectory(repositoryPath); !isDirectory {
		return nil, fmt.Errorf("The supplied item repository path %q is not a directory.", repositoryPath)
	}

	// create the file provider
	fileProvider, err := newFileProvider(logger, repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create the item provider because the file provider could not be created. Error: %s", err.Error())
	}

	return &itemProvider{
		logger:         logger,
		repositoryPath: repositoryPath,
		fileProvider:   fileProvider,
	}, nil
}

type itemProvider struct {
	logger         logger.Logger
	repositoryPath string

	fileProvider *fileProvider
}

func (itemProvider *itemProvider) GetItemFromDirectory(itemDirectory string) (item *dataaccess.Item, err error) {

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

func (itemProvider *itemProvider) getChildItemsFromDirectory(itemDirectory string) (childItems []*dataaccess.Item) {

	childItems = make([]*dataaccess.Item, 0)

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

func (itemProvider *itemProvider) newItemFromFile(itemDirectory, filePath string) (item *dataaccess.Item, err error) {
	// route
	route, err := route.NewFromItemPath(itemProvider.repositoryPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a physical item from route %q", route)

	// content
	contentProvider := newFileContentProvider(filePath, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []*dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item, err = dataaccess.NewPhysicalItem(route, contentProvider, files, childs)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	return item, nil
}

func (itemProvider *itemProvider) newVirtualItem(itemDirectory string) (item *dataaccess.Item, err error) {

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a virtual item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []*dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item, err = dataaccess.NewVirtualItem(route, contentProvider, files, childs)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	return item, nil
}

func (itemProvider *itemProvider) newFileCollectionItem(itemDirectory string) (item *dataaccess.Item, err error) {

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a file collection item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s

files: [Attachments](files)`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := itemDirectory
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// create the item
	item, err = dataaccess.NewFileCollectionItem(route, contentProvider, files)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	return item, nil
}
