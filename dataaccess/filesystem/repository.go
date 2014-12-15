// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"path/filepath"
	"runtime"
	"time"
)

type UpdateCallback func(route.Route)

type event struct {
	Item  *dataaccess.Item
	Error error
}

func newRepositoryEvent(item *dataaccess.Item, err error) event {
	return event{
		Item:  item,
		Error: err,
	}
}

type Repository struct {
	logger    logger.Logger
	hash      string
	directory string

	itemProvider *itemProvider

	// Indizes
	items          []*dataaccess.Item
	itemByRoute    map[string]*dataaccess.Item
	itemByHash     map[string]*dataaccess.Item
	itemHasChanged map[string]bool

	// Updates
	reindexNotificationChannels []chan bool
	routesWithSubscribers       map[string]route.Route
	onUpdateCallbacks           []UpdateCallback
}

func NewRepository(logger logger.Logger, directory string, reindexIntervalInSeconds int) (*Repository, error) {

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

	// create the repository
	repository := &Repository{
		logger:    logger,
		directory: directory,

		itemProvider: itemProvider,

		// Indizes
		items:          make([]*dataaccess.Item, 0),
		itemByRoute:    make(map[string]*dataaccess.Item),
		itemByHash:     make(map[string]*dataaccess.Item),
		itemHasChanged: make(map[string]bool),

		// Updates
		routesWithSubscribers:       make(map[string]route.Route),
		onUpdateCallbacks:           make([]UpdateCallback, 0),
		reindexNotificationChannels: make([]chan bool, 0),
	}

	// index the repository
	repository.init()

	// scheduled reindex
	repository.reindex(reindexIntervalInSeconds)

	return repository, nil
}

func (repository *Repository) Path() string {
	return repository.directory
}

func (repository *Repository) Items() []*dataaccess.Item {
	return repository.items
}

func (repository *Repository) Item(route route.Route) *dataaccess.Item {
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

	newItemsChannel := make(chan event, 1)

	go func() {
		repository.discoverItems(repository.directory, newItemsChannel)
		close(newItemsChannel)
	}()

	newItems := make([]*dataaccess.Item, 0)
	newItemByRoute := make(map[string]*dataaccess.Item)
	newItemByHash := make(map[string]*dataaccess.Item)
	newItemHasChanged := make(map[string]bool)

	for event := range newItemsChannel {

		err := event.Error
		if err != nil {
			repository.logger.Warn(err.Error())
		}

		item := event.Item
		if item == nil {
			repository.logger.Warn("The even contained an empty item but no error.")
			continue
		}

		// determine the item hash
		hash, err := item.Hash()
		if err != nil {
			repository.logger.Warn("Could not determine the hash for item %q. Error: %s", item.Route(), err.Error())
			continue
		}

		// check if the hash changed
		hasChanged := false
		if _, itemByHashWasFound := repository.itemByHash[hash]; !itemByHashWasFound {
			hasChanged = true // it's a new item
		}

		repository.logger.Debug("Adding item %q", item)

		newItems = append(newItems, item)
		newItemByRoute[item.Route().Value()] = item
		newItemHasChanged[item.Route().Value()] = hasChanged
		newItemByHash[hash] = item
	}

	repository.items = newItems
	repository.itemByRoute = newItemByRoute
	repository.itemHasChanged = newItemHasChanged
	repository.itemByHash = newItemByHash

	// inform subscribers about updates
	repository.notifySubscribers()

	// send out after reindex notifications
	repository.sendAfterReindexUpdates()
}

func (repository *Repository) AfterReindex(notificationChannel chan bool) {
	repository.reindexNotificationChannels = append(repository.reindexNotificationChannels, notificationChannel)
}

func (repository *Repository) OnUpdate(callback func(route.Route)) {
	repository.onUpdateCallbacks = append(repository.onUpdateCallbacks, callback)
}

func (repository *Repository) StartWatching(r route.Route) {
	repository.routesWithSubscribers[route.ToKey(r)] = r

	// todo: start a watcher
}

func (repository *Repository) StopWatching(r route.Route) {
	key := route.ToKey(r)
	delete(repository.routesWithSubscribers, key)

	// todo: stop the watcher started earlier
}

func (repository *Repository) sendAfterReindexUpdates() {
	for _, updateChannel := range repository.reindexNotificationChannels {
		updateChannel <- true
	}
}

func (repository *Repository) notifySubscribers() {
	for _, route := range repository.routesWithSubscribers {

		hasChanged, exists := repository.itemHasChanged[route.Value()]
		if !exists {
			// don't notify if the item does no longer exists
			continue
		}

		if !hasChanged {
			// don't notify if the item has not changed
			continue
		}

		repository.logger.Info("Item %q has changed.", route)

		for _, onUpdateCallback := range repository.onUpdateCallbacks {
			go onUpdateCallback(route)
		}

	}
}

// Start the fulltext search indexing process
func (repository *Repository) reindex(intervalInSeconds int) {

	if intervalInSeconds <= 1 {
		repository.logger.Debug("Reindexing is disabled.")
		return
	}

	go func() {
		sleepInterval := time.Second * time.Duration(intervalInSeconds)
		for {

			repository.logger.Info("Number of go routines: %d", runtime.NumGoroutine())
			repository.logger.Info("Reindexing")
			repository.init()

			time.Sleep(sleepInterval)
		}
	}()
}

// Create a new Item for the specified path.
func (repository *Repository) discoverItems(itemDirectory string, targetChannel chan event) {

	// create the item
	item, err := repository.itemProvider.GetItemFromDirectory(itemDirectory)

	// send the item to the target channel
	targetChannel <- newRepositoryEvent(item, err)

	// abort if the item cannot have childs
	if !item.CanHaveChilds() {
		return
	}

	// recurse for child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		repository.discoverItems(childItemDirectory, targetChannel)
	}
}
