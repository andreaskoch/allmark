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
	"time"
)

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
	items []*dataaccess.Item

	// Updates
	reindexNotificationChannels []chan bool
	routesWithSubscribers       map[string]route.Route
	onUpdateCallback            func(route.Route)
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
		items: make([]*dataaccess.Item, 0),

		// Updates
		routesWithSubscribers:       make(map[string]route.Route),
		onUpdateCallback:            func(r route.Route) {},
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

	for event := range newItemsChannel {

		err := event.Error
		if err != nil {
			repository.logger.Warn(err.Error())
		}

		item := event.Item
		if item == nil {
			continue
		}

		repository.logger.Debug("Adding item %q", item)
		newItems = append(newItems, item)
	}

	repository.items = newItems

	// inform subscribers about updates
	repository.notifySubscribers()

	// send out after reindex notifications
	repository.sendAfterReindexUpdates()
}

func (repository *Repository) AfterReindex() chan bool {
	notificationChannel := make(chan bool, 1)
	repository.reindexNotificationChannels = append(repository.reindexNotificationChannels, notificationChannel)
	return notificationChannel
}

func (repository *Repository) OnUpdate(callback func(route.Route)) {
	repository.onUpdateCallback = callback
}

func (repository *Repository) StartWatching(r route.Route) {
	repository.routesWithSubscribers[route.ToKey(r)] = r
}

func (repository *Repository) StopWatching(r route.Route) {
	key := route.ToKey(r)
	delete(repository.routesWithSubscribers, key)
}

func (repository *Repository) sendAfterReindexUpdates() {
	for _, updateChannel := range repository.reindexNotificationChannels {
		go func() {
			updateChannel <- true
		}()
	}
}

func (repository *Repository) notifySubscribers() {
	// for _, route := range repository.routesWithSubscribers {
	// 	go repository.onUpdateCallback(route)
	// }
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
