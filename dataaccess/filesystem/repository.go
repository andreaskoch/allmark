// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem/index"
	"github.com/andreaskoch/go-fswatch"
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

	onUpdateCallback func(route.Route)
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

	if logger.Level() == loglevel.Debug {

		// enable the debug mode for the filesystem watcher
		debugMessages := fswatch.EnableDebug()
		go func() {
			for message := range debugMessages {
				logger.Debug("fs-watch: %s", message)
			}
		}()
	}

	// create the repository
	repository := &Repository{
		logger:    logger,
		directory: directory,
		hash:      hash,

		itemProvider: itemProvider,
		index:        index.New(logger),

		onUpdateCallback: func(r route.Route) {},
	}

	// index the repository
	repository.init()

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

func (repository *Repository) OnUpdate(callback func(route.Route)) {
	repository.onUpdateCallback = callback
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

func (repository *Repository) StartWatching(route route.Route) {
	repository.itemProvider.StartWatching(route)
}

func (repository *Repository) StopWatching(route route.Route) {
	repository.itemProvider.StopWatching(route)
}

// Initialize the repository - scan all folders and update the index.
func (repository *Repository) init() {

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

		repository.logger.Info("Adding item %q", item)
		repository.index.Add(item)
	}

	// scheduled reindex of the fulltext index
	repository.startFullTextSearch()

	// start watching for updates
	repository.startWatching()
}

func (repository *Repository) notifySubscribers(route route.Route) {
	go repository.onUpdateCallback(route)
}

// Start the fulltext search indexing process
func (repository *Repository) startFullTextSearch() {

	repository.itemSearch = dataaccess.NewItemSearch(repository.logger, repository)

	go func() {
		sleepInterval := time.Minute * 3
		for {
			repository.logger.Info("Refreshing the search index.")
			repository.itemSearch.Update()

			time.Sleep(sleepInterval)
		}
	}()
}

// Start watching for updates
func (repository *Repository) startWatching() {
	updateChannel := repository.itemProvider.Updates()

	go func() {

		for {
			select {
			case newItemEvent := <-updateChannel.New:
				{
					item := newItemEvent.Item
					err := newItemEvent.Error

					if err != nil {

						repository.logger.Warn("New Item. Error: %s", err)

					} else if item != nil {

						// Add item to index
						repository.logger.Info("New item %q", item)
						repository.index.Add(item)

						// Send out updates
						go repository.onUpdateCallback(item.Route())

					}
				}

			case changedItemEvent := <-updateChannel.Changed:
				{
					item := changedItemEvent.Item
					err := changedItemEvent.Error

					if err != nil {

						repository.logger.Warn("Changed Item. Error: %s", err)

					} else if item != nil {

						// Add item to index
						repository.logger.Info("Changed item %q", item)
						repository.index.Add(item)

						// Send out updates
						go repository.onUpdateCallback(item.Route())

					}
				}

			case movedItemEvent := <-updateChannel.Moved:
				{
					item := movedItemEvent.Item
					err := movedItemEvent.Error

					if err != nil {

						repository.logger.Warn("Moved Item. Error: %s", err)

					} else if item != nil {

						repository.logger.Info("Moved item %q", item)
						repository.index.Remove(item.Route())

					}
				}

			}
		}

	}()
}

// Create a new Item for the specified path.
func (repository *Repository) discoverItems(itemDirectory string, targetChannel chan event) {

	// create the item
	item, recurse, err := repository.itemProvider.GetItemFromDirectory(itemDirectory)

	// send the item to the target channel
	targetChannel <- newRepositoryEvent(item, err)

	// recurse for child items
	if recurse {
		childItemDirectories := getChildDirectories(itemDirectory)
		for _, childItemDirectory := range childItemDirectories {
			repository.discoverItems(childItemDirectory, targetChannel)
		}
	}
}
