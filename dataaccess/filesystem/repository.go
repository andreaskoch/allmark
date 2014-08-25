// Copyright 2014 Andreas Koch. All rights reserved.
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
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Repository struct {
	logger    logger.Logger
	hash      string
	directory string
	watcher   *watcherFactory

	index      *Index
	itemSearch *dataaccess.ItemSearch

	updateHub *UpdateHub

	newItem     chan *dataaccess.RepositoryEvent // new items which are discovered after the first index has been built
	changedItem chan *dataaccess.RepositoryEvent // items with changed content
	movedItem   chan *dataaccess.RepositoryEvent // items which moved

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

	// hash provider: use the directory name for the hash (for now)
	directoryName := strings.ToLower(filepath.Base(directory))
	hash, err := getStringHash(directoryName)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a hash for the repository with the name %q. Error: %s", directoryName, err)
	}

	// enable debug mode
	// debugMessages := fswatch.EnableDebug()
	// go func() {
	// 	for message := range debugMessages {
	// 		logger.Debug("fs-watch: %s", message)
	// 	}
	// }()

	// create the repository
	repository := &Repository{
		logger:    logger,
		directory: directory,
		hash:      hash,

		index:     newIndex(logger),
		updateHub: newUpdateHub(logger),
		watcher:   newWatcherFactory(logger),

		// item channels
		newItem:     make(chan *dataaccess.RepositoryEvent, 1),
		changedItem: make(chan *dataaccess.RepositoryEvent, 1),
		movedItem:   make(chan *dataaccess.RepositoryEvent, 1),

		onUpdateCallback: func(r route.Route) {},
	}

	// index the repository
	repository.startItemDiscovery()

	// scheduled reindex of the fulltext index
	repository.startFullTextSearch()

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
	repository.updateHub.StartWatching(route)
}

func (repository *Repository) StopWatching(route route.Route) {
	repository.updateHub.StopWatching(route)
}

// Start indexing the the repository
func (repository *Repository) startItemDiscovery() {

	go func() {

		for {
			select {
			case newItemEvent := <-repository.newItem:
				{
					item := newItemEvent.Item
					err := newItemEvent.Error

					if err != nil {

						repository.logger.Warn("New Item. Error: %s", err)

					} else if item != nil {

						// Add item to index
						repository.logger.Info("New item %q", item)
						repository.index.Add(item)

						// Inform subscribers about change
						repository.notifySubscribers(item.Route())

					}
				}

			case changedItemEvent := <-repository.changedItem:
				{
					item := changedItemEvent.Item
					err := changedItemEvent.Error

					if err != nil {

						repository.logger.Warn("Changed Item. Error: %s", err)

					} else if item != nil {

						// Add item to index
						repository.logger.Info("Changed item %q", item)
						repository.index.Add(item)

						// Inform subscribers about change
						repository.notifySubscribers(item.Route())

					}
				}

			case movedItemEvent := <-repository.movedItem:
				{
					item := movedItemEvent.Item
					err := movedItemEvent.Error

					if err != nil {

						repository.logger.Warn("Moved Item. Error: %s", err)

					} else if item != nil {

						repository.logger.Info("Moved item %q", item)
						repository.index.Remove(item.Route())

					} else {

						repository.logger.Warn("Empty repository event in the moved item channel.")

					}
				}

			}
		}

	}()

	repository.discoverItems(repository.directory, repository.newItem)
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
	item, recurse := repository.getItemFromDirectory(repository.Path(), itemDirectory)

	// send the item to the target channel
	targetChannel <- dataaccess.NewEvent(item, nil)

	// recurse for child items
	if recurse {
		childItemDirectories := getChildDirectories(itemDirectory)
		for _, childItemDirectory := range childItemDirectories {
			repository.discoverItems(childItemDirectory, targetChannel)
		}
	}
}

func (repository *Repository) getItemFromDirectory(repositoryPath, itemDirectory string) (item *dataaccess.Item, recurse bool) {

	// physical item from markdown file
	if found, markdownFilePath := findMarkdownFileInDirectory(itemDirectory); found {

		// create an item from the markdown file
		return repository.newItemFromFile(repositoryPath, itemDirectory, markdownFilePath)

	}

	// virtual item
	if directoryContainsItems(itemDirectory, 3) {
		return repository.newVirtualItem(repositoryPath, itemDirectory)
	}

	// file collection item
	return repository.newFileCollectionItem(repositoryPath, itemDirectory)
}

func (repository *Repository) newItemFromFile(repositoryPath, itemDirectory, filePath string) (item *dataaccess.Item, recurse bool) {
	// route
	route, err := route.NewFromItemPath(repositoryPath, filePath)
	if err != nil {
		repository.logger.Error("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
		return
	}

	repository.logger.Info("Creating a physical item from route %q", route)

	// content provider
	checkIntervalInSeconds := 2
	contentProvider := newFileContentProvider(filePath, route, checkIntervalInSeconds)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: Sub-Directory Watcher
	repository.updateHub.Attach(route, "sub-directory-watcher", func() fswatch.Watcher {
		return repository.watcher.SubDirectories(itemDirectory, 2, func(change *fswatch.FolderChange) {

			// remove the parent item since we cannot easily determine which child has gone away
			go func() {
				repository.movedItem <- dataaccess.NewEvent(item, nil)
			}()

			// recreate this item
			repository.discoverItems(itemDirectory, repository.newItem)
		})
	})

	// Update-Hub: File-Directory Watcher
	repository.updateHub.Attach(route, "file-directory-watcher", func() fswatch.Watcher {
		return repository.watcher.AllFiles(filesDirectory, 2, func(change *fswatch.FolderChange) {

			// update file list
			repository.logger.Debug("Updating the file list for item %q", item.String())
			newFiles := getFiles(repositoryPath, itemDirectory, filesDirectory)
			item.SetFiles(newFiles)

			// update the parent item
			go func() {
				repository.changedItem <- dataaccess.NewEvent(item, nil)
			}()
		})
	})

	// Update-Hub: Markdown-File Watcher
	repository.updateHub.Attach(route, "markdown-file-watcher", func() fswatch.Watcher {

		modifiedCallback := func() {
			repository.logger.Debug("Item %q changed.", item)
			go func() {
				repository.changedItem <- dataaccess.NewEvent(item, nil)
			}()
		}

		movedCallback := func() {
			repository.logger.Debug("Item %q moved.", item)
			go func() {
				repository.movedItem <- dataaccess.NewEvent(item, nil)
			}()
		}

		return repository.watcher.File(filePath, 2, modifiedCallback, movedCallback)
	})

	return item, true
}

func (repository *Repository) newVirtualItem(repositoryPath, itemDirectory string) (item *dataaccess.Item, recurse bool) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

	// route
	route, err := route.NewFromItemDirectory(repositoryPath, itemDirectory)
	if err != nil {
		repository.logger.Error("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
		return
	}

	repository.logger.Info("Creating a virtual item from route %q", route)

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: Sub-Directory Watcher
	repository.updateHub.Attach(route, "sub-directory-watcher", func() fswatch.Watcher {
		return repository.watcher.SubDirectories(itemDirectory, 2, func(change *fswatch.FolderChange) {

			// remove the parent item since we cannot easily determine which child has gone away
			go func() {
				repository.movedItem <- dataaccess.NewEvent(item, nil)
			}()

			// recreate this item
			repository.discoverItems(itemDirectory, repository.newItem)
		})
	})

	// Update-Hub: Type-Change Watcher
	repository.updateHub.Attach(route, "type-change-watcher", func() fswatch.Watcher {
		return repository.watcher.Directory(itemDirectory, 2, func(change *fswatch.FolderChange) {

			for _, newFile := range change.New() {
				if !isMarkdownFile(newFile) {
					continue
				}

				// remove the parent item since we cannot easily determine which child has gone away
				go func() {
					repository.movedItem <- dataaccess.NewEvent(item, nil)
				}()

				// recreate this item
				repository.discoverItems(itemDirectory, repository.newItem)

				break

			}
		})
	})

	// Update-Hub: File-Directory Watcher
	repository.updateHub.Attach(route, "file-directory-watcher", func() fswatch.Watcher {
		return repository.watcher.AllFiles(filesDirectory, 2, func(change *fswatch.FolderChange) {

			// update file list
			repository.logger.Debug("Updating the file list for item %q", item.String())
			newFiles := getFiles(repositoryPath, itemDirectory, filesDirectory)
			item.SetFiles(newFiles)

			// update the parent item
			go func() {
				repository.changedItem <- dataaccess.NewEvent(item, nil)
			}()
		})
	})

	return item, true
}

func (repository *Repository) newFileCollectionItem(repositoryPath, itemDirectory string) (item *dataaccess.Item, recurse bool) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

	// route
	route, err := route.NewFromItemDirectory(repositoryPath, itemDirectory)
	if err != nil {
		repository.logger.Error("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
		return
	}

	repository.logger.Info("Creating a file collection item from route %q", route)

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := itemDirectory
	files := getFiles(repositoryPath, itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: File-Change Watcher
	repository.updateHub.Attach(route, "file-change-watcher", func() fswatch.Watcher {
		return repository.watcher.AllFiles(itemDirectory, 2, func(change *fswatch.FolderChange) {

			if directoryContainsItems(itemDirectory, 1) {
				// type change

				// remove the parent item since we cannot easily determine which child has gone away
				go func() {
					repository.movedItem <- dataaccess.NewEvent(item, nil)
				}()

				// recreate this item
				repository.discoverItems(itemDirectory, repository.newItem)

			} else {

				// update files
				newFiles := getFiles(repositoryPath, itemDirectory, filesDirectory)
				item.SetFiles(newFiles)

				// update the parent item
				go func() {
					repository.changedItem <- dataaccess.NewEvent(item, nil)
				}()
			}
		})
	})

	return item, false
}

// Check if the specified directory contains an item within the range of the given max depth.
func directoryContainsItems(directory string, maxdepth int) bool {

	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		childDirectory := filepath.Join(directory, entry.Name())

		if entry.IsDir() {
			if isReservedDirectory(childDirectory) {
				continue
			}

			if maxdepth > 0 {

				// recurse
				if directoryContainsItems(childDirectory, maxdepth-1) {
					return true
				}
			}

			continue
		}

		if isMarkdownFile(childDirectory) {
			return true
		}

		continue
	}

	return false
}
