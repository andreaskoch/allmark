// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/util/fsutil"
	"github.com/elWyatt/allmark/dataaccess"
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
		repository.logger.Info("Reindexing: On")
		repository.reindex(config.Indexing.IntervalInSeconds)
	} else {
		repository.logger.Info("Reindexing: Off")
	}

	// live reload
	if config.LiveReload.Enabled {
		repository.logger.Info("Live Reload: On")
	} else {
		repository.logger.Info("Live Reload: Off")
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
	routes := make([]route.Route, 0)

	for _, item := range repository.index.GetAllItems() {
		routes = append(routes, item.Route())
	}

	return routes
}

// Subscribe registers the supplied updates channel in the repository.
// All updates (new, modified or deleted items) in the repository will be passed down this channel.
func (repository *Repository) Subscribe(updates chan dataaccess.Update) {
	repository.updateSubscribers = append(repository.updateSubscribers, updates)
}

// StartWatching starts the watcher for the item with the given route.
func (repository *Repository) StartWatching(route route.Route) {

	if !repository.livereloadIsEnabled {
		repository.logger.Info("Live reload: Off")
		return
	}

	item := repository.Item(route)
	if item == nil {
		repository.logger.Warn("Cannot start watching. Item %q was not found.", route.String())
		return
	}

	itemRoute := item.Route()

	fileSystemItem := item.(*Item)
	updates, err := repository.watcher.Start(itemRoute, fileSystemItem.WatcherPaths())
	if err != nil {
		return
	}

	itemDirectory := fileSystemItem.Directory()

	go func() {
		for repository.watcher.IsRunning(route) {
			select {
			case <-updates:

				repository.logger.Info("Received an update for route %q. Rescanning directory %q.", itemRoute, itemDirectory)

				// update the index
				oldIndex := repository.index
				limitDepth := true
				maxDepth := 2
				repository.updateIndex(oldIndex, itemRoute, itemDirectory, limitDepth, maxDepth)

			}
		}
	}()
}

// StopWatching stops the watcher for the item with the given route.
func (repository *Repository) StopWatching(route route.Route) {
	repository.watcher.Stop(route)
}

// Initialize the repository - scan all folders and update the index.
func (repository *Repository) init() {

	var oldIndex *Index
	if repository.index != nil {
		repository.logger.Debug("Re-initializing the repository index.")
		oldIndex = repository.index
	} else {
		repository.logger.Debug("Initializing the repository index.")
		oldIndex = newIndex()
	}

	limitDepth := false // we want to index all items
	maxDepth := 0

	repository.updateIndex(oldIndex, route.New(), repository.directory, limitDepth, maxDepth)
}

// createIndexFromDirectory scans the supplied directory and creates an index from it.
// If limitMaxDepth is set to true maxDepth defines the max depth of the scan and of the resulting index.
func (repository *Repository) createIndexFromDirectory(directory string, limitMaxDepth bool, maxDepth int) *Index {

	repository.logger.Debug("Scanning directory %q", directory)

	index := newIndex()

	// update the cloned index
	for _, newItem := range repository.getItemsFromDirectory(directory, limitMaxDepth, maxDepth) {

		if _, err := index.Add(newItem); err != nil {
			repository.logger.Error("Cannot add item %q to index: Error: %s", newItem.String, err.Error())
		}

	}

	return index
}

// getItemsFromDirectory scan the supplied directory for items and returns the list if items found.
// If limitMaxDepth is set to true maxDepth defines the max depth of the scan.
func (repository *Repository) getItemsFromDirectory(itemDirectory string, limitDepth bool, maxDepth int) (items []dataaccess.Item) {

	items = make([]dataaccess.Item, 0)

	// create the item
	item, err := repository.itemProvider.GetItemFromDirectory(itemDirectory)
	if err != nil {
		repository.logger.Error("Could not create an item from folder %q. Error: %s", itemDirectory, err.Error())
		return
	}

	// append the item
	items = append(items, item)

	// abort if the item cannot have children
	if !item.CanHaveChildren() {
		return
	}

	if limitDepth {

		// abort if the max depth level has been reached
		if maxDepth == 0 {
			return items
		}

		// count down the max depth
		maxDepth = maxDepth - 1

	}

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		childItems := repository.getItemsFromDirectory(childItemDirectory, limitDepth, maxDepth)
		items = append(items, childItems...)
	}

	return
}

// reindex starts the scheduled reindexing process.
func (repository *Repository) reindex(intervalInSeconds int) {

	if intervalInSeconds <= 1 {
		repository.logger.Info("Reindexing: Off")
		return
	}

	go func() {
		sleepInterval := time.Second * time.Duration(intervalInSeconds)

		for {

			// wait for the next turn
			time.Sleep(sleepInterval)

			repository.logger.Debug("Number of go routines: %d", runtime.NumGoroutine())
			repository.logger.Info("Reindexing")

			// index
			repository.init()
		}
	}()
}

// sendUpdate send an update to all subscribers.
func (repository *Repository) sendUpdate(update dataaccess.Update) {
	if update.IsEmpty() {
		repository.logger.Debug("sendUpdate(%s): Nothing to send.", update.String())
		return
	}

	repository.logger.Debug("sendUpdate(%s): Notifying all subscribers", update.String())
	for _, updateSubscriber := range repository.updateSubscribers {
		updateSubscriber <- update
	}
}

// updateIndex takes the supplied oldIndex and an updates it with items it found in the specified directory.
// If limitMaxDepth is set to true maxDepth defines the max depth of the scan.
func (repository *Repository) updateIndex(oldIndex *Index, itemRoute route.Route, itemDirectory string, limitDepth bool, maxDepth int) {

	// get the old sub index
	subIndexOld := oldIndex.GetSubIndex(itemRoute, limitDepth, maxDepth)

	// get the new sub index
	subIndexNew := repository.createIndexFromDirectory(itemDirectory, limitDepth, maxDepth)

	repository.logger.Debug("------- Sub Indexes for %q ---------------", itemDirectory)
	repository.logger.Debug("Sub index (old):\n%s", subIndexOld.String())
	repository.logger.Debug("Sub index (new):\n%s", subIndexNew.String())

	// determine the diff between the old and new sub indexes
	newItems, modifiedItems, deletedItems := repository.diffIndexes(subIndexOld, subIndexNew)

	repository.logger.Debug("------- Difference ---------------")
	repository.logger.Debug("New: %v", len(newItems))
	repository.logger.Debug("Modified: %v", len(modifiedItems))
	repository.logger.Debug("Deleted: %v", len(deletedItems))

	// prepare the new index
	newIndex := oldIndex.Copy()

	// remove deleted items
	for _, deletedItem := range deletedItems {
		newIndex.Remove(deletedItem.Route())
	}

	// add new items
	for _, newItem := range newItems {
		newIndex.Add(newItem)
	}

	repository.logger.Debug("------- Full Indexes ---------------")
	repository.logger.Debug("Old Index:\n%s", oldIndex.String())
	repository.logger.Debug("New Index:\n%s", newIndex.String())

	// assign the new index
	repository.index = newIndex

	// send out updates
	changedItems := dataaccess.NewUpdate(itemsToRoutes(newItems), itemsToRoutes(modifiedItems), itemsToRoutes(deletedItems))
	repository.sendUpdate(changedItems)
}

// diffIndexes calculates the differences between the specified old and new indexes.
func (repository *Repository) diffIndexes(oldIndex, newIndex *Index) (newItems, modifiedItems, deletedItems []dataaccess.Item) {

	// new or modified
	for _, newItem := range newIndex.GetAllItems() {

		oldItem, existsInOldIndex := oldIndex.IsMatch(newItem.Route())
		if !existsInOldIndex {
			// it's new
			repository.logger.Debug("%q was not found in\n%s", newItem.Route(), oldIndex.String())
			repository.logger.Debug("%q: new", newItem.Route())

			newItems = append(newItems, newItem)
			continue
		}

		// check if it has changed
		// determine the hash of the new item
		newItemHash, err := newItem.Hash()
		if err != nil {
			repository.logger.Error("Skipping item %q because the hash of the new item cannot be determined", newItem.Route())
			continue
		}

		// determine the hash of the old item
		oldItemHash := oldItem.LastHash()
		if oldItemHash == newItemHash {
			continue
		}

		repository.logger.Debug("%q: modified", newItem.Route())
		modifiedItems = append(modifiedItems, newItem)
	}

	// deleted
	for _, oldItem := range oldIndex.GetAllItems() {
		_, existsInNewIndex := newIndex.IsMatch(oldItem.Route())
		if existsInNewIndex {
			// never mind
			continue
		}

		// it's deleted
		repository.logger.Debug("%q: deleted", oldItem.Route())
		deletedItems = append(deletedItems, oldItem)
	}

	return

}

// itemsToRoutes returns the list of routes for a given listen if items.
func itemsToRoutes(items []dataaccess.Item) []route.Route {
	var routes []route.Route
	for _, item := range items {
		routes = append(routes, item.Route())
	}
	return routes
}
