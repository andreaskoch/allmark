// Copyright 2013 Andreas Koch. All rights reserved.
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
	"strings"
)

type Repository struct {
	logger    logger.Logger
	hash      string
	directory string
}

func NewRepository(logger logger.Logger, directory string) (*Repository, error) {

	// check if path exists
	if !fsutil.PathExists(directory) {
		return nil, fmt.Errorf("The path %q does not exist.", directory)
	}

	// check if the supplied path is a file
	if isDirectory, _ := fsutil.IsDirectory(directory); !isDirectory {
		directory = filepath.Dir(directory)
	}

	// abort if the supplied path is a reserved directory
	if isReservedDirectory(directory) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be a root.", directory)
	}

	// hash provider: use the directory name for the hash (for now)
	directoryName := strings.ToLower(filepath.Base(directory))
	hash, err := getStringHash(directoryName)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a hash for the repository with the name %q. Error: %s", directoryName, err)
	}

	return &Repository{
		logger:    logger,
		directory: directory,
		hash:      hash,
	}, nil
}

func (repository *Repository) GetItems() (itemEvents chan *dataaccess.RepositoryEvent) {

	itemEvents = make(chan *dataaccess.RepositoryEvent, 1)

	go func() {

		// repository directory item
		indexItems(repository, repository.directory, itemEvents)

		// close the channel. All items have been indexed
		close(itemEvents)
	}()

	return itemEvents
}

func (repository *Repository) Id() string {
	return repository.hash
}

func (repository *Repository) Path() string {
	return repository.directory
}

// Create a new Item for the specified path.
func indexItems(repository *Repository, itemPath string, itemEvents chan *dataaccess.RepositoryEvent) {

	// abort if path does not exist
	if !fsutil.PathExists(itemPath) {
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q does not exist.", itemPath))
		return
	}

	// abort if path is reserved
	if isReservedDirectory(itemPath) {
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemPath))
		return
	}

	// make sure the item path points to a markdown file
	isVirtualItem := false
	itemDirectory := filepath.Dir(itemPath)
	if isDirectory, _ := fsutil.IsDirectory(itemPath); isDirectory {

		// search for a markdown file in the directory
		if found, filepath := findMarkdownFileInDirectory(itemPath); found {

			itemDirectory = itemPath
			itemPath = filepath

		} else {

			// virtual item
			isVirtualItem = true
			itemDirectory = itemPath

		}

	} else if !isMarkdownFile(itemPath) {

		// the supplied item path does not point to a markdown file
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("%q is not a markdown file.", itemPath))
		return
	}

	// create a new item
	if !isVirtualItem {
		// route
		route, err := route.NewFromItemPath(repository.Path(), itemPath)
		if err != nil {
			itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemPath, err))
		}

		// content provider
		contentProvider := newContentProvider(itemPath, route)

		// create the file index
		filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
		files := getFiles(repository, itemDirectory, filesDirectory)

		// create the item
		item, err := dataaccess.NewItem(route, contentProvider, files)

		itemEvents <- dataaccess.NewEvent(item, err)
	}

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		indexItems(repository, childItemDirectory, itemEvents)
	}
}
