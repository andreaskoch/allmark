// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/util/fsutil"
	"github.com/elWyatt/allmark/dataaccess"
	"fmt"
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

	route := itemProvider.GetRouteFromFilePath(filePath)
	itemProvider.logger.Debug("Creating a physical item from route %q", route)

	// content
	contentProvider, contentProviderError := newFileContentProvider(filePath, route)
	if contentProviderError != nil {
		return nil, contentProviderError
	}

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// children
	children := func() []dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item := newPhysicalItem(
		route,
		contentProvider,
		files,
		children,
		itemDirectory,
		[]watcherPather{
			watcherFilePath{filePath},
			watcherDirectoryPath{itemDirectory, false},
			watcherDirectoryPath{filesDirectory, true},
		},
	)
	return item, nil
}

func (itemProvider *itemProvider) newVirtualItem(itemDirectory string) (dataaccess.Item, error) {

	route := itemProvider.GetRouteFromDirectory(itemDirectory)
	itemProvider.logger.Debug("Creating a virtual item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider, contentProviderError := newTextContentProvider(content, route)
	if contentProviderError != nil {
		return nil, contentProviderError
	}

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// children
	children := func() []dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item := newVirtualItem(
		route,
		contentProvider,
		files,
		children,
		itemDirectory,
		[]watcherPather{
			watcherDirectoryPath{itemDirectory, false},
		})

	return item, nil
}

func (itemProvider *itemProvider) newFileCollectionItem(itemDirectory string) (dataaccess.Item, error) {

	route := itemProvider.GetRouteFromDirectory(itemDirectory)
	itemProvider.logger.Debug("Creating a file collection item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s

files: [Attachments](/)`, title)
	contentProvider, contentProviderError := newTextContentProvider(content, route)
	if contentProviderError != nil {
		return nil, contentProviderError
	}

	// files
	filesDirectory := itemDirectory
	files := func() []dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// create the item
	item := newFileCollectionItem(
		route,
		contentProvider,
		files,
		itemDirectory,
		[]watcherPather{
			watcherDirectoryPath{itemDirectory, true},
		},
	)

	return item, nil
}

// GetRouteFromDirectory creates a route from the given directory path.
func (itemProvider *itemProvider) GetRouteFromDirectory(directory string) route.Route {
	return route.NewFromItemDirectory(itemProvider.repositoryPath, directory)
}

// GetRouteFromFilePath creates a route from the given file path.
func (itemProvider *itemProvider) GetRouteFromFilePath(filepath string) route.Route {
	return route.NewFromItemPath(itemProvider.repositoryPath, filepath)
}
