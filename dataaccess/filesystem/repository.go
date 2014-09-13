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
	"github.com/andreaskoch/allmark2/dataaccess/filesystem/index"
	"path/filepath"
	"strings"
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

	index      *index.Index
	itemSearch *dataaccess.ItemSearch

	routesWithSubscribers map[string]route.Route
	onUpdateCallback      func(route.Route)
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

	// hash provider: use the directory name for the hash (for now)
	directoryName := strings.ToLower(filepath.Base(directory))
	hash, err := getStringHash(directoryName)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a hash for the repository with the name %q. Error: %s", directoryName, err)
	}

	// create the repository
	repository := &Repository{
		logger:    logger,
		directory: directory,
		hash:      hash,

		itemProvider: itemProvider,

		routesWithSubscribers: make(map[string]route.Route),
		onUpdateCallback:      func(r route.Route) {},
	}

	// index the repository
	repository.init()

	// scheduled reindex
	repository.reindex(reindexIntervalInSeconds)

	return repository, nil
}

func (repository *Repository) String() string {
	return repository.index.String()
}

func (repository *Repository) Root() *dataaccess.Item {
	return repository.index.Root()
}

func (repository *Repository) Items() []*dataaccess.Item {
	return repository.index.Items()
}

func (repository *Repository) Item(route route.Route) (*dataaccess.Item, bool) {
	return repository.index.IsMatch(route)
}

func (repository *Repository) File(route route.Route) (*dataaccess.File, bool) {
	return repository.index.IsFileMatch(route)
}

func (repository *Repository) Parent(route route.Route) *dataaccess.Item {
	return repository.index.GetParent(route)
}

func (repository *Repository) Childs(route route.Route) []*dataaccess.Item {
	return repository.index.GetDirectChilds(route)
}

func (repository *Repository) AllChilds(route route.Route) []*dataaccess.Item {
	return repository.index.GetAllChilds(route, func(item *dataaccess.Item) bool {
		return true
	})
}

func (repository *Repository) AllMatchingChilds(route route.Route, matchExpression func(item *dataaccess.Item) bool) []*dataaccess.Item {
	return repository.index.GetAllChilds(route, matchExpression)
}

func (repository *Repository) Search(keywords string, maxiumNumberOfResults int) (searchResults []dataaccess.SearchResult) {

	if repository.itemSearch == nil {
		repository.logger.Warn("The fulltext index has not been initialized (Keyword: %q).", keywords)
		return
	}

	return repository.itemSearch.Search(keywords, maxiumNumberOfResults)
}

func (repository *Repository) Id() string {
	return repository.hash
}

func (repository *Repository) Path() string {
	return repository.directory
}

func (repository *Repository) Size() int {
	return repository.index.Size()
}

// Initialize the repository - scan all folders and update the index.
func (repository *Repository) init() {

	index := index.New(repository.logger)
	newItems := make(chan event, 1)

	go func() {
		repository.discoverItems(repository.directory, newItems)
		close(newItems)
	}()

	for event := range newItems {

		err := event.Error
		if err != nil {
			repository.logger.Warn(err.Error())
		}

		item := event.Item
		if item == nil {
			continue
		}

		repository.logger.Debug("Adding item %q", item)
		index.Add(item)
	}

	// assign the new index
	repository.index = index

	// update search index
	repository.logger.Debug("Refreshing the search index.")

	// cleanup
	previousItemSearch := repository.itemSearch
	if previousItemSearch != nil {
		defer func() {
			previousItemSearch.Destroy()
		}()
	}

	repository.itemSearch = dataaccess.NewItemSearch(repository.logger, repository)

	// inform subscribers about updates
	repository.notifySubscribers()
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
