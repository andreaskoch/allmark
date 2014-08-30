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
	"github.com/andreaskoch/allmark2/dataaccess/filesystem/updates"
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"path/filepath"
)

type itemUpdateChannel struct {
	Moved   chan event
	Changed chan event
	New     chan event
}

func newItemUpdateChannel() *itemUpdateChannel {
	return &itemUpdateChannel{
		Moved:   make(chan event, 1),
		Changed: make(chan event, 1),
		New:     make(chan event, 1),
	}
}

func newItemProvider(logger logger.Logger, repositoryPath string) (*itemProvider, error) {

	// abort if repoistory path does not exist
	if !fsutil.PathExists(repositoryPath) {
		return nil, fmt.Errorf("The repository path %q does not exist.", repositoryPath)
	}

	// abort if the supplied repository path is not a directory
	if isDirectory, _ := fsutil.IsDirectory(repositoryPath); !isDirectory {
		return nil, fmt.Errorf("The supplied item repository path %q is not a directory.", repositoryPath)
	}

	// create the file provider
	fileProvider, err := newFileProvider(logger, repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create the item provider because the file provider could not be created. Error: %s", err.Error())
	}

	return &itemProvider{
		logger:         logger,
		repositoryPath: repositoryPath,
		fileProvider:   fileProvider,

		updateChannel: newItemUpdateChannel(),
		updateHub:     updates.NewHub(logger),
		watcher:       newWatcherFactory(logger),
	}, nil
}

type itemProvider struct {
	logger         logger.Logger
	repositoryPath string

	fileProvider *fileProvider

	updateChannel *itemUpdateChannel
	updateHub     *updates.Hub
	watcher       *watcherFactory
}

func (itemProvider *itemProvider) Updates() *itemUpdateChannel {
	return itemProvider.updateChannel
}

func (itemProvider *itemProvider) StartWatching(route route.Route) {
	itemProvider.updateHub.StartWatching(route)
}

func (itemProvider *itemProvider) StopWatching(route route.Route) {
	itemProvider.updateHub.StopWatching(route)
}

func (itemProvider *itemProvider) GetItemFromDirectory(itemDirectory string) (item *dataaccess.Item, recurse bool, err error) {

	// abort if path does not exist
	if !fsutil.PathExists(itemDirectory) {
		return nil, false, fmt.Errorf("The path %q does not exist.", itemDirectory)
	}

	// make sure the item directory points to a folder not a file
	if isDirectory, _ := fsutil.IsDirectory(itemDirectory); !isDirectory {
		itemProvider.logger.Warn("The supplied item directory path %q is not a directory using the parent instead.", itemDirectory)
		itemDirectory = filepath.Dir(itemDirectory)
	}

	// abort if path is reserved
	if isReservedDirectory(itemDirectory) {
		return nil, false, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemDirectory)
	}

	// physical item from markdown file
	if found, markdownFilePath := findMarkdownFileInDirectory(itemDirectory); found {

		// create an item from the markdown file
		return itemProvider.newItemFromFile(itemDirectory, markdownFilePath)

	}

	// virtual item
	if directoryContainsItems(itemDirectory, 3) {
		return itemProvider.newVirtualItem(itemDirectory)
	}

	// file collection item
	return itemProvider.newFileCollectionItem(itemDirectory)
}

func (itemProvider *itemProvider) newItemFromFile(itemDirectory, filePath string) (item *dataaccess.Item, recurse bool, err error) {
	// route
	route, err := route.NewFromItemPath(itemProvider.repositoryPath, filePath)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a physical item from route %q", route)

	// content provider
	checkIntervalInSeconds := 2
	contentProvider := newFileContentProvider(filePath, route, checkIntervalInSeconds)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	itemProvider.updateHub.Attach(route, "sub-directory-watcher", itemProvider.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	itemProvider.updateHub.Attach(route, "file-directory-watcher", itemProvider.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	// Update-Hub: Markdown-File Watcher
	itemProvider.updateHub.Attach(route, "markdown-file-watcher", itemProvider.fileWatcher(item, filePath))

	return item, true, nil
}

func (itemProvider *itemProvider) newVirtualItem(itemDirectory string) (item *dataaccess.Item, recurse bool, err error) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a virtual item from route %q", route)

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	itemProvider.updateHub.Attach(route, "sub-directory-watcher", itemProvider.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: Type-Change Watcher
	itemProvider.updateHub.Attach(route, "type-change-watcher", itemProvider.newMarkdownFileWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	itemProvider.updateHub.Attach(route, "file-directory-watcher", itemProvider.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	return item, true, nil
}

func (itemProvider *itemProvider) newFileCollectionItem(itemDirectory string) (item *dataaccess.Item, recurse bool, err error) {

	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a file collection item from route %q", route)

	// content provider
	contentProvider := newTextContentProvider(content, route)

	// create the file index
	filesDirectory := itemDirectory
	files := itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		return nil, false, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: File-Change Watcher
	itemProvider.updateHub.Attach(route, "file-change-watcher", itemProvider.newMarkdownFileWatcher(item, itemDirectory))

	return item, true, nil
}

func (itemProvider *itemProvider) onStartTriggerFunc(item *dataaccess.Item, itemDirectory, filesDirectory string) func() {
	return func() {

		// update files
		newFiles := itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
		item.SetFiles(newFiles)

		go func() {
			itemProvider.updateChannel.Changed <- newRepositoryEvent(item, nil)
		}()

	}
}

func (itemProvider *itemProvider) fileDirectoryWatcher(item *dataaccess.Item, itemDirectory, filesDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return itemProvider.watcher.AllFiles(filesDirectory, 2, func(change *fswatch.FolderChange) {

			// update file list
			itemProvider.logger.Debug("Updating the file list for item %q", item.String())
			newFiles := itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
			item.SetFiles(newFiles)

			go func() {
				itemProvider.updateChannel.Changed <- newRepositoryEvent(item, nil)
			}()
		})
	}

}

func (itemProvider *itemProvider) subDirectoryWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return itemProvider.watcher.SubDirectories(itemDirectory, 2, func(change *fswatch.FolderChange) {
			// itemProvider.discoverItems(itemDirectory, itemProvider.newItem)
		})
	}

}

func (itemProvider *itemProvider) newMarkdownFileWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return itemProvider.watcher.Directory(itemDirectory, 2, func(change *fswatch.FolderChange) {

			// check if one of the files is a markdown file
			oneOfTheNewFilesIsAMarkdownFile := false
			for _, newFile := range change.New() {
				if isMarkdownFile(newFile) {
					oneOfTheNewFilesIsAMarkdownFile = true
					break
				}
			}

			// no change if there is no markdown file
			if !oneOfTheNewFilesIsAMarkdownFile {
				return
			}

			// reindex this item
			// itemProvider.discoverItems(itemDirectory, itemProvider.changedItem)
		})
	}
}

func (itemProvider *itemProvider) fileWatcher(item *dataaccess.Item, filePath string) func() fswatch.Watcher {

	return func() fswatch.Watcher {

		modifiedCallback := func() {
			itemProvider.logger.Debug("Item %q changed.", item)
			go func() {
				itemProvider.updateChannel.Changed <- newRepositoryEvent(item, nil)
			}()
		}

		movedCallback := func() {
			itemProvider.logger.Debug("Item %q moved.", item)
			go func() {
				itemProvider.updateChannel.Moved <- newRepositoryEvent(item, nil)
			}()
		}

		return itemProvider.watcher.File(filePath, 2, modifiedCallback, movedCallback)
	}
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
