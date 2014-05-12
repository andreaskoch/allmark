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
		indexItems(repository.Path(), repository.Path(), itemEvents)

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
func indexItems(repositoryPath, itemPath string, itemEvents chan *dataaccess.RepositoryEvent) {

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

	// make sure the item directory points to a folder not a file
	itemDirectory := itemPath
	if isDirectory, _ := fsutil.IsDirectory(itemPath); !isDirectory {
		itemDirectory = filepath.Dir(itemPath)
	}

	// search for a markdown file in the directory
	if found, markdownFilePath := findMarkdownFileInDirectory(itemPath); found {

		// create an item from the markdown file
		item, err := newItemFromFile(repositoryPath, itemDirectory, markdownFilePath)
		itemEvents <- dataaccess.NewEvent(item, err)

	} else {

		// create a virtual item
		title := filepath.Base(itemDirectory)

		content := fmt.Sprintf(`# %s`, title)

		item, err := newItemFromText(repositoryPath, itemDirectory, content)
		itemEvents <- dataaccess.NewEvent(item, err)

	}

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		indexItems(repositoryPath, childItemDirectory, itemEvents)
	}
}

func newItemFromFile(repositoryPath, itemDirectory, filePath string) (*dataaccess.Item, error) {
	// route
	route, err := route.NewFromItemPath(repositoryPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", filePath, err)
	}

	// content provider
	contentProvider := newFileContentProvider(filePath, route)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	return dataaccess.NewItem(route, contentProvider, files)
}

func newItemFromText(repositoryPath, itemDirectory, content string) (*dataaccess.Item, error) {
	// route
	route, err := route.NewFromItemDirectory(repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	return dataaccess.NewItem(route, contentProvider, files)
}
