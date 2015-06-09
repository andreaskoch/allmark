// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/dataaccess"
)

type Repository struct {
	logger    logger.Logger
	directory string

	itemProvider *itemProvider

	index *Index

	// Update Subscription
	watcher           *filesystemWatcher
	updateSubscribers []chan dataaccess.Update

	// live reload
	livereloadIsEnabled bool
}

func NewRepository(logger logger.Logger, directory string, config config.Config) (*Repository, error) {

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
		index: newIndex(),

		// Update Subscription
		watcher:           newFilesystemWatcher(logger),
		updateSubscribers: updateSubscribers,

		livereloadIsEnabled: config.LiveReload.Enabled,
	}

	// index the repository
	repository.init()

	// scheduled reindex
	if config.Indexing.Enabled {
		repository.reindex(config.Indexing.IntervalInSeconds)
	} else {
		repository.logger.Info("Reindexing is disabled.")
	}

	return repository, nil
}

func (repository *Repository) Path() string {
	return repository.directory
}

func (repository *Repository) Items() []dataaccess.Item {
	return repository.index.GetAllItems()
}

func (repository *Repository) Item(route route.Route) dataaccess.Item {
	item, isMatch := repository.index.IsMatch(route)
	if !isMatch {
		return nil
	}
	return item
}

func (repository *Repository) Routes() []route.Route {
	return getRoutesFromIndex(repository.index)
}

func getRoutesFromIndex(index *Index) []route.Route {
	routes := make([]route.Route, 0)

	for _, item := range index.GetAllItems() {
		routes = append(routes, item.Route())
	}

	return routes
}

// Initialize the repository - scan all folders and update the index.
func (repository *Repository) init() {
	newIndex, _ := repository.rescan(repository.directory, false, 0)
	repository.index = newIndex
}

func (repository *Repository) rescan(directory string, limitMaxDepth bool, maxDepth int) (*Index, dataaccess.Update) {

	repository.logger.Debug("Rescanning", directory)

	// notification listssrc/allmark.io/modules/web/handlers/update.go
	newItemRoutes := make([]route.Route, 0)
	modifiedItemRoutes := make([]route.Route, 0)
	deletedItemRoutes := make([]route.Route, 0)

	index := repository.index.Copy()

	// scan the repository directory for new items
	for _, newItem := range repository.getItemsFromDirectory(directory, limitMaxDepth, maxDepth) {

		// check if the item is new or modified
		existingItem := repository.Item(newItem.Route())
		isNewItem := existingItem == nil
		if isNewItem {

			// the route was not found in the index it must be a new item
			newItemRoutes = append(newItemRoutes, newItem.Route())
			repository.logger.Debug("Item %q is new", newItem.Route())

		} else {

			// compare hashes
			hash := newItem.LastHash()
			existingItemHash := existingItem.LastHash()

			if existingItemHash != hash {

				// the hash has changed
				modifiedItemRoutes = append(modifiedItemRoutes, newItem.Route())
				repository.logger.Debug("Item %q is changed", newItem.Route())

			} else {
				repository.logger.Debug("Item %q is not modified", newItem.Route())
			}

		}

		index.Add(newItem)

	}

	// find deleted items
	for _, oldItemRoute := range repository.Routes() {
		if _, exists := index.IsMatch(oldItemRoute); exists {
			continue
		}

		deletedItemRoutes = append(deletedItemRoutes, oldItemRoute)

		// remove the item from the existing index
		index.Remove(oldItemRoute)
	}

	// send update to subscribers
	return index, dataaccess.NewUpdate(newItemRoutes, modifiedItemRoutes, deletedItemRoutes)

}

// Create a new Item for the specified path.
func (repository *Repository) getItemsFromDirectory(itemDirectory string, limitDepth bool, maxDepth int) (items []dataaccess.Item) {

	items = make([]dataaccess.Item, 0)

	if limitDepth {

		// abort if the max depth level has been reached
		if maxDepth == 0 {
			return items
		}

		// count down the max depth
		maxDepth = maxDepth - 1

	}

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
		childItems := repository.getItemsFromDirectory(childItemDirectory, limitDepth, maxDepth)
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
		repository.logger.Debug("sendUpdate: Empty update")
		return
	}

	repository.logger.Debug("sendUpdate: %s", update.String())
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

				repository.logger.Info("Recieved an update for route %q. Reindexing.", route.String())

				newIndex, updates := repository.rescan(fileSystemItem.Directory(), true, 1)

				repository.index = newIndex
				repository.sendUpdate(updates)

			}
		}
	}()
}

func (repository *Repository) StopWatching(route route.Route) {
	repository.watcher.Stop(route)
}
