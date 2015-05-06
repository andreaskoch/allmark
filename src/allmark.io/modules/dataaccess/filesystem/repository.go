// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/dataaccess"
)

type Repository struct {
	logger    logger.Logger
	hash      string
	directory string

	itemProvider *itemProvider

	// Indizes
	items       []dataaccess.Item
	itemByRoute map[string]dataaccess.Item
	itemByHash  map[string]dataaccess.Item

	// Update Subscription
	watcher           *filesystemWatcher
	updateSubscribers []chan dataaccess.Update

	// live reload
	livereloadIsEnabled bool
}

func NewRepository(logger logger.Logger, directory string, reindexIntervalInSeconds int, reindex, livereload bool) (*Repository, error) {

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

	itemProvider, err := newItemProvider(logger, directory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create the repository because the item provider could not be created. Error: %s", err.Error())
	}

	// create an update channel
	updateSubscribers := make([]chan dataaccess.Update, 0)

	// create the repository
	repository := &Repository{
		logger:    logger,
		directory: directory,

		itemProvider: itemProvider,

		// Indizes
		items:       make([]dataaccess.Item, 0),
		itemByRoute: make(map[string]dataaccess.Item),
		itemByHash:  make(map[string]dataaccess.Item),

		// Update Subscription
		watcher:             newFilesystemWatcher(logger),
		updateSubscribers:   updateSubscribers,
		livereloadIsEnabled: livereload,
	}

	// index the repository
	repository.init()

	// scheduled reindex
	if reindex {
		repository.reindex(reindexIntervalInSeconds)
	} else {
		repository.logger.Info("Reindexing is disabled.")
	}

	return repository, nil
}

func (repository *Repository) Path() string {
	return repository.directory
}

func (repository *Repository) Items() []dataaccess.Item {
	return repository.items
}

func (repository *Repository) Item(route route.Route) dataaccess.Item {
	return repository.itemByRoute[route.Value()]
}

func (repository *Repository) Routes() []route.Route {

	routes := make([]route.Route, 0, len(repository.items))
	for _, item := range repository.items {
		routes = append(routes, item.Route())
	}

	return routes
}

// Initialize the repository - scan all folders and update the index.
func (repository *Repository) init() {

	// notification lists
	newItemRoutes := make([]route.Route, 0)
	modifiedItemRoutes := make([]route.Route, 0)
	deletedItemRoutes := make([]route.Route, 0)

	// create new indices
	items := make([]dataaccess.Item, 0)
	itemByRoute := make(map[string]dataaccess.Item)
	itemByHash := make(map[string]dataaccess.Item)

	// scan the repository directory for new items
	for _, item := range repository.getItemsFromDirectory(repository.directory) {

		// determine the item hash
		hash, err := item.Hash()
		if err != nil {
			repository.logger.Warn("Could not determine the hash for item %q. Error: %s", item.Route(), err.Error())
			continue
		}

		// check if the item is new or modified
		existingItem := repository.Item(item.Route())
		isNewItem := existingItem == nil
		if isNewItem {

			// the route was not found in the index it must be a new item
			newItemRoutes = append(newItemRoutes, item.Route())

		} else {

			// check if the hash is already in the index
			if _, itemHashIsAlreadyInTheIndex := repository.itemByHash[hash]; itemHashIsAlreadyInTheIndex == false {

				// the item has changed the hash was not found in the index
				modifiedItemRoutes = append(modifiedItemRoutes, item.Route())
			}

		}

		// store the item in the indizes
		items = append(items, item)
		itemByRoute[item.Route().Value()] = item
		itemByHash[hash] = item
	}

	// find deleted items
	for _, oldItem := range repository.items {
		oldItemRoute := oldItem.Route()
		if _, oldItemExistsInNewIndex := itemByRoute[oldItemRoute.Value()]; oldItemExistsInNewIndex {
			continue
		}

		deletedItemRoutes = append(deletedItemRoutes, oldItemRoute)
	}

	// override the existing values
	repository.items = items
	repository.itemByRoute = itemByRoute
	repository.itemByHash = itemByHash

	// send update to subscribers
	repository.sendUpdate(dataaccess.NewUpdate(newItemRoutes, modifiedItemRoutes, deletedItemRoutes))
}

// Create a new Item for the specified path.
func (repository *Repository) getItemsFromDirectory(itemDirectory string) (items []dataaccess.Item) {

	items = make([]dataaccess.Item, 0)

	// create the item
	item, err := repository.itemProvider.GetItemFromDirectory(itemDirectory)
	if err != nil {
		repository.logger.Error("Could not create an item from folder %q", itemDirectory)
		return
	}

	// append the item
	items = append(items, item)

	// abort if the item cannot have childs
	if !item.CanHaveChilds() {
		return
	}

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		childItems := repository.getItemsFromDirectory(childItemDirectory)
		items = append(items, childItems...)
	}

	return
}

// Start the fulltext search indexing process
func (repository *Repository) reindex(intervalInSeconds int) {

	if intervalInSeconds <= 1 {
		repository.logger.Info("Reindexing: Off")
		return
	}

	repository.logger.Info("Reindexing: On")

	go func() {
		sleepInterval := time.Second * time.Duration(intervalInSeconds)

		for {

			repository.logger.Debug("Number of go routines: %d", runtime.NumGoroutine())
			repository.logger.Info("Reindexing")

			// index
			repository.init()

			// wait for the next turn
			time.Sleep(sleepInterval)
		}
	}()
}

func (repository *Repository) Subscribe(updates chan dataaccess.Update) {
	repository.updateSubscribers = append(repository.updateSubscribers, updates)
}

// Send an update down the repository update channel
func (repository *Repository) sendUpdate(update dataaccess.Update) {
	if update.IsEmpty() {
		return
	}

	for _, updateSubscriber := range repository.updateSubscribers {
		updateSubscriber <- update
	}
}

func (repository *Repository) StartWatching(route route.Route) {

	if !repository.livereloadIsEnabled {
		repository.logger.Info("Live reload is disabled.")
		return
	}

	item := repository.Item(route)
	if item == nil {
		repository.logger.Warn("Cannot start watching. Item %q was not found.", route.String())
		return
	}

	fileSystemItem := item.(*Item)
	updates, err := repository.watcher.Start(fileSystemItem.Route(), fileSystemItem.WatcherPaths())
	if err != nil {
		return
	}

	go func() {
		for repository.watcher.IsRunning(route) {
			select {
			case <-updates:

				repository.logger.Info("Sending out an update for route %q.", route.String())
				repository.sendUpdate(dataaccess.NewModifiedItemUpdate(route))

			}
		}
	}()
}

func (repository *Repository) StopWatching(route route.Route) {
	repository.watcher.Stop(route)
}
