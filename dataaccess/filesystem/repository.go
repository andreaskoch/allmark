// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Repository struct {
	logger    logger.Logger
	hash      string
	directory string

	newItem     chan *dataaccess.RepositoryEvent // new items which are discovered after the first index has been built
	changedItem chan *dataaccess.RepositoryEvent // items with changed content
	movedItem   chan *dataaccess.RepositoryEvent // items which moved
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

		newItem:     make(chan *dataaccess.RepositoryEvent, 1),
		changedItem: make(chan *dataaccess.RepositoryEvent, 1),
		movedItem:   make(chan *dataaccess.RepositoryEvent, 1),
	}, nil
}

func (repository *Repository) InitialItems() chan *dataaccess.RepositoryEvent {

	// open the channel
	startupItem := make(chan *dataaccess.RepositoryEvent, 1)

	go func() {

		// repository directory item
		repository.discoverItems(repository.Path(), startupItem)

		// close the channel. All items have been indexed
		close(startupItem)
	}()

	return startupItem
}

func (repository *Repository) NewItems() chan *dataaccess.RepositoryEvent {
	return repository.newItem
}

func (repository *Repository) ChangedItems() chan *dataaccess.RepositoryEvent {
	return repository.changedItem
}

func (repository *Repository) MovedItems() chan *dataaccess.RepositoryEvent {
	return repository.movedItem
}

func (repository *Repository) Id() string {
	return repository.hash
}

func (repository *Repository) Path() string {
	return repository.directory
}

// Create a new Item for the specified path.
func (repository *Repository) discoverItems(itemPath string, targetChannel chan *dataaccess.RepositoryEvent) {

	// abort if path does not exist
	if !fsutil.PathExists(itemPath) {
		targetChannel <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q does not exist.", itemPath))
		return
	}

	// abort if path is reserved
	if isReservedDirectory(itemPath) {
		targetChannel <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemPath))
		return
	}

	// make sure the item directory points to a folder not a file
	itemDirectory := itemPath
	if isDirectory, _ := fsutil.IsDirectory(itemPath); !isDirectory {
		itemDirectory = filepath.Dir(itemPath)
	}

	// create the item
	item, err := getItemFromDirectory(repository.Path(), itemDirectory)
	if err != nil {
		repository.logger.Warn("Unable to create item from directory %q. Error: %s", itemDirectory, err.Error())
		return // abort
	}

	// send the item to the target channel
	targetChannel <- dataaccess.NewEvent(item, nil)

	// attach content change listener
	repository.attachContentListener(item)

	// attach directory listener
	repository.attachDirectoryListener(itemDirectory)

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		repository.discoverItems(childItemDirectory, targetChannel)
	}
}

func (repository *Repository) attachDirectoryListener(itemDirectory string) {

	// look for changes in the item directory
	go func() {
		var skipFunc = func(path string) bool {
			isReserved := isReservedDirectory(path)
			return isReserved
		}

		folderWatcher := fswatch.NewFolderWatcher(itemDirectory, false, skipFunc).Start()

		for folderWatcher.IsRunning() {

			select {
			case <-folderWatcher.Change:
				repository.logger.Info("Folder %q changed", itemDirectory)
				repository.discoverItems(itemDirectory, repository.newItem)
			}

		}
	}()
}

func (repository *Repository) attachContentListener(item *dataaccess.Item) {

	// watch for changes
	go func() {
		contentChangeChannel := item.ChangeEvent()
	ChannelLoop:
		for changeEvent := range contentChangeChannel {

			switch changeEvent {

			case content.TypeChanged:
				{
					repository.logger.Info("Item %q changed.", item)
					repository.changedItem <- dataaccess.NewEvent(item, nil)
				}

			case content.TypeMoved:
				{
					repository.logger.Info("Item %q moved.", item)
					repository.movedItem <- dataaccess.NewEvent(item, nil)

					break ChannelLoop
				}

			}

		}

		repository.logger.Debug("Exiting content listener for item %q.", item)
	}()
}

func getItemFromDirectory(repositoryPath, itemDirectory string) (*dataaccess.Item, error) {

	// physical item from markdown file
	if found, markdownFilePath := findMarkdownFileInDirectory(itemDirectory); found {

		// create an item from the markdown file
		return newItemFromFile(repositoryPath, itemDirectory, markdownFilePath)

	}

	// virtual item
	if directoryDoesNotContainsItems(itemDirectory) {
		return newVirtualItem(repositoryPath, itemDirectory)
	}

	// file collection item
	return newFileCollectionItem(repositoryPath, itemDirectory)
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

func newVirtualItem(repositoryPath, itemDirectory string) (*dataaccess.Item, error) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

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

func newFileCollectionItem(repositoryPath, itemDirectory string) (*dataaccess.Item, error) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

	// route
	route, err := route.NewFromItemDirectory(repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := itemDirectory
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	return dataaccess.NewItem(route, contentProvider, files)
}

func directoryDoesNotContainsItems(directory string) bool {
	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		childDirectory := filepath.Join(directory, entry.Name())

		if entry.IsDir() {
			if isReservedDirectory(childDirectory) {
				return true
			} else {
				return directoryDoesNotContainsItems(childDirectory)
			}
		} else if isMarkdownFile(childDirectory) {
			return true
		} else {
			continue
		}
	}

	return false
}
