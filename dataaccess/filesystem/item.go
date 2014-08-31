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
	"path/filepath"
)

type itemUpdateChannel struct {
	Moved   chan route.Route
	Changed chan route.Route
	New     chan event
}

func newItemUpdateChannel() *itemUpdateChannel {
	return &itemUpdateChannel{
		Moved:   make(chan route.Route, 1),
		Changed: make(chan route.Route, 1),
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

func (itemProvider *itemProvider) GetItemFromDirectory(itemDirectory string) (item *dataaccess.Item, err error) {

	// abort if path does not exist
	if !fsutil.PathExists(itemDirectory) {
		return nil, fmt.Errorf("The path %q does not exist.", itemDirectory)
	}

	// make sure the item directory points to a folder not a file
	if isDirectory, _ := fsutil.IsDirectory(itemDirectory); !isDirectory {
		itemProvider.logger.Warn("The supplied item directory path %q is not a directory using the parent instead.", itemDirectory)
		itemDirectory = filepath.Dir(itemDirectory)
	}

	// abort if path is reserved
	if isReservedDirectory(itemDirectory) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemDirectory)
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

func (itemProvider *itemProvider) getChildItemsFromDirectory(itemDirectory string) (childItems []*dataaccess.Item) {

	childItems = make([]*dataaccess.Item, 0)

	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		child, err := itemProvider.GetItemFromDirectory(childItemDirectory)
		if err != nil {
			itemProvider.logger.Warn("Cannot create item from directory. Error: %s", err.Error())
		}

		childItems = append(childItems, child)
	}

	return childItems
}

func (itemProvider *itemProvider) newItemFromFile(itemDirectory, filePath string) (item *dataaccess.Item, err error) {
	// route
	route, err := route.NewFromItemPath(itemProvider.repositoryPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a physical item from route %q", route)

	// content
	contentProvider := newFileContentProvider(filePath, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []*dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item, err = dataaccess.NewPhysicalItem(route, contentProvider, files, childs)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	itemProvider.updateHub.Attach(route, "sub-directory-watcher", itemProvider.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	itemProvider.updateHub.Attach(route, "file-directory-watcher", itemProvider.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	// Update-Hub: Markdown-File Watcher
	itemProvider.updateHub.Attach(route, "markdown-file-watcher", itemProvider.fileWatcher(item, filePath))

	return item, nil
}

func (itemProvider *itemProvider) newVirtualItem(itemDirectory string) (item *dataaccess.Item, err error) {

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a virtual item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// childs
	childs := func() []*dataaccess.Item {
		return itemProvider.getChildItemsFromDirectory(itemDirectory)
	}

	// create the item
	item, err = dataaccess.NewVirtualItem(route, contentProvider, files, childs)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	itemProvider.updateHub.Attach(route, "sub-directory-watcher", itemProvider.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: Type-Change Watcher
	itemProvider.updateHub.Attach(route, "type-change-watcher", itemProvider.newMarkdownFileWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	itemProvider.updateHub.Attach(route, "file-directory-watcher", itemProvider.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	return item, nil
}

func (itemProvider *itemProvider) newFileCollectionItem(itemDirectory string) (item *dataaccess.Item, err error) {

	// route
	route, err := route.NewFromItemDirectory(itemProvider.repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create an Item for the path %q. Error: %s", itemDirectory, err)
	}

	itemProvider.logger.Info("Creating a file collection item from route %q", route)

	// content
	title := filepath.Base(itemDirectory)
	content := fmt.Sprintf(`# %s`, title)
	contentProvider := newTextContentProvider(content, route)

	// files
	filesDirectory := itemDirectory
	files := func() []*dataaccess.File {
		return itemProvider.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
	}

	// create the item
	item, err = dataaccess.NewFileCollectionItem(route, contentProvider, files)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Item %q. Error: %s", route, err)
	}

	// Update-Hub: OnStart Trigger
	itemProvider.updateHub.RegisterOnStartTrigger(route, itemProvider.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: File-Change Watcher
	itemProvider.updateHub.Attach(route, "file-change-watcher", itemProvider.newMarkdownFileWatcher(item, itemDirectory))

	return item, nil
}

func (itemProvider *itemProvider) onStartTriggerFunc(item *dataaccess.Item, itemDirectory, filesDirectory string) func() {

	itemRoute := item.Route()
	previousChilds := item.GetChilds()

	return func() {

		go func() {
			newChilds, removedChilds := item.GetChildItemChanges(previousChilds)
			for _, childRoute := range removedChilds {
				itemProvider.logger.Debug("Route %q has moved.", childRoute)
				itemProvider.updateChannel.Moved <- childRoute
			}

			for _, child := range newChilds {
				itemProvider.logger.Debug("New child %q .", child)
				itemProvider.updateChannel.New <- newRepositoryEvent(child, nil)
			}

			// update the previous childs list for the next tiem
			previousChilds = item.GetChilds()

			itemProvider.updateChannel.Changed <- itemRoute
		}()

	}
}

func (itemProvider *itemProvider) fileDirectoryWatcher(item *dataaccess.Item, itemDirectory, filesDirectory string) func() fswatch.Watcher {

	itemRoute := item.Route()

	return func() fswatch.Watcher {
		return itemProvider.watcher.AllFiles(filesDirectory, 2, func(change *fswatch.FolderChange) {

			go func() {
				itemProvider.updateChannel.Changed <- itemRoute
			}()
		})
	}

}

func (itemProvider *itemProvider) subDirectoryWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	itemRoute := item.Route()

	return func() fswatch.Watcher {
		return itemProvider.watcher.SubDirectories(itemDirectory, 2, func(change *fswatch.FolderChange) {

			// new items
			for _, newFolder := range change.New() {

				itemProvider.logger.Debug("The folder %q has been detected as a new child of %q.", newFolder, itemDirectory)
				newItem, err := itemProvider.GetItemFromDirectory(newFolder)
				if err != nil {
					itemProvider.logger.Warn(err.Error())
					continue
				}

				itemProvider.updateChannel.New <- newRepositoryEvent(newItem, nil)
			}

			// moved items
			for _, movedFolder := range change.Moved() {
				itemProvider.logger.Debug("Folder %q has moved.", movedFolder)

				movedItemRoute, err := route.NewFromItemDirectory(itemProvider.repositoryPath, movedFolder)

				itemProvider.logger.Debug("Route of the moved folder is %q.", movedItemRoute)

				if err != nil {
					itemProvider.logger.Warn(err.Error())
					continue
				}

				itemProvider.updateChannel.Moved <- movedItemRoute
			}

			// update the parent
			if len(change.New()) > 0 || len(change.Moved()) > 0 {

				go func() {
					itemProvider.updateChannel.Changed <- itemRoute
				}()

			}

		})
	}

}

func (itemProvider *itemProvider) newMarkdownFileWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	itemRoute := item.Route()

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

			go func() {
				itemProvider.updateChannel.Changed <- itemRoute
			}()
		})
	}
}

func (itemProvider *itemProvider) fileWatcher(item *dataaccess.Item, filePath string) func() fswatch.Watcher {

	itemRoute := item.Route()

	return func() fswatch.Watcher {

		modifiedCallback := func() {
			itemProvider.logger.Debug("Item %q changed.", itemRoute)
			go func() {
				itemProvider.updateChannel.Changed <- itemRoute
			}()
		}

		movedCallback := func() {
			itemProvider.logger.Debug("Item %q moved.", itemRoute)
			go func() {
				itemProvider.updateChannel.Moved <- itemRoute
			}()
		}

		return itemProvider.watcher.File(filePath, 2, modifiedCallback, movedCallback)
	}
}
