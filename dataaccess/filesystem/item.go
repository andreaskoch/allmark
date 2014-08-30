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
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"path/filepath"
)

func newItemProvider(logger logger.Logger, repositoryPath string) (*itemProvider, error) {

	// abort if repoistory path does not exist
	if !fsutil.PathExists(itemDirectory) {
		return nil, fmt.Errorf("The repository path %q does not exist.", repositoryPath), false
	}

	// abort if the supplied repository path is not a directory
	if isDirectory, _ := fsutil.IsDirectory(itemDirectory); !isDirectory {
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
	}, nil
}

type itemProvider struct {
	logger         logger.Logger
	repositoryPath string

	fileProvider *fileProvider
}

func (itemProvider *itemProvider) GetItemFromDirectory(itemDirectory string) (item *dataaccess.Item, err error, recurse bool) {

	// abort if path does not exist
	if !fsutil.PathExists(itemDirectory) {
		return nil, fmt.Errorf("The path %q does not exist.", itemDirectory), false
	}

	// make sure the item directory points to a folder not a file
	if isDirectory, _ := fsutil.IsDirectory(itemDirectory); !isDirectory {
		itemProvider.logger.Warn("The supplied item directory path %q is not a directory using the parent instead.", itemDirectory)
		itemDirectory = filepath.Dir(itemDirectory)
	}

	// abort if path is reserved
	if isReservedDirectory(itemDirectory) {
		targetChannel <- newRepositoryEvent(nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemDirectory))
		return
	}

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

func (itemProvider *itemProvider) newItemFromFile(repositoryPath, itemDirectory, filePath string) (item *dataaccess.Item, recurse bool) {
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
	files := repository.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: OnStart Trigger
	repository.updateHub.RegisterOnStartTrigger(route, repository.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	repository.updateHub.Attach(route, "sub-directory-watcher", repository.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	repository.updateHub.Attach(route, "file-directory-watcher", repository.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	// Update-Hub: Markdown-File Watcher
	repository.updateHub.Attach(route, "markdown-file-watcher", repository.fileWatcher(item, filePath))

	return item, true
}

func (itemProvider *itemProvider) newVirtualItem(repositoryPath, itemDirectory string) (item *dataaccess.Item, recurse bool) {

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
	files := repository.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: OnStart Trigger
	repository.updateHub.RegisterOnStartTrigger(route, repository.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: Sub-Directory Watcher
	repository.updateHub.Attach(route, "sub-directory-watcher", repository.subDirectoryWatcher(item, itemDirectory))

	// Update-Hub: Type-Change Watcher
	repository.updateHub.Attach(route, "type-change-watcher", repository.newMarkdownFileWatcher(item, itemDirectory))

	// Update-Hub: File-Directory Watcher
	repository.updateHub.Attach(route, "file-directory-watcher", repository.fileDirectoryWatcher(item, itemDirectory, filesDirectory))

	return item, true
}

func (itemProvider *itemProvider) newFileCollectionItem(repositoryPath, itemDirectory string) (item *dataaccess.Item, recurse bool) {

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
	files := repository.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)

	// create the item
	item, err = dataaccess.NewItem(route, contentProvider, files)
	if err != nil {
		repository.logger.Error("Cannot create Item %q. Error: %s", route, err)
		return
	}

	// Update-Hub: OnStart Trigger
	repository.updateHub.RegisterOnStartTrigger(route, repository.onStartTriggerFunc(item, itemDirectory, filesDirectory))

	// Update-Hub: File-Change Watcher
	repository.updateHub.Attach(route, "file-change-watcher", repository.newMarkdownFileWatcher(item, itemDirectory))

	return item, false
}

func (itemProvider *itemProvider) onStartTriggerFunc(item *dataaccess.Item, itemDirectory, filesDirectory string) func() {
	return func() {

		// update files
		newFiles := repository.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
		item.SetFiles(newFiles)

		go func() {
			repository.changedItem <- newRepositoryEvent(item, nil)
		}()

	}
}

func (itemProvider *itemProvider) fileDirectoryWatcher(item *dataaccess.Item, itemDirectory, filesDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return repository.watcher.AllFiles(filesDirectory, 2, func(change *fswatch.FolderChange) {

			// update file list
			repository.logger.Debug("Updating the file list for item %q", item.String())
			newFiles := repository.fileProvider.GetFilesFromDirectory(itemDirectory, filesDirectory)
			item.SetFiles(newFiles)

			go func() {
				repository.changedItem <- newRepositoryEvent(item, nil)
			}()
		})
	}

}

func (itemProvider *itemProvider) subDirectoryWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return repository.watcher.SubDirectories(itemDirectory, 2, func(change *fswatch.FolderChange) {
			repository.discoverItems(itemDirectory, repository.newItem)
		})
	}

}

func (itemProvider *itemProvider) newMarkdownFileWatcher(item *dataaccess.Item, itemDirectory string) func() fswatch.Watcher {

	return func() fswatch.Watcher {
		return repository.watcher.Directory(itemDirectory, 2, func(change *fswatch.FolderChange) {

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
			repository.discoverItems(itemDirectory, repository.changedItem)
		})
	}
}

func (itemProvider *itemProvider) fileWatcher(item *dataaccess.Item, filePath string) func() fswatch.Watcher {

	return func() fswatch.Watcher {

		modifiedCallback := func() {
			repository.logger.Debug("Item %q changed.", item)
			go func() {
				repository.changedItem <- newRepositoryEvent(item, nil)
			}()
		}

		movedCallback := func() {
			repository.logger.Debug("Item %q moved.", item)
			go func() {
				repository.movedItem <- newRepositoryEvent(item, nil)
			}()
		}

		return repository.watcher.File(filePath, 2, modifiedCallback, movedCallback)
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
