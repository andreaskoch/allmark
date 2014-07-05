// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/go-fswatch"
)

var (
	watchers map[string]fswatch.Watcher
)

func init() {
	watchers = make(map[string]fswatch.Watcher)
}

func newWatcherFactory(logger logger.Logger) *watcherFactory {
	return &watcherFactory{
		logger: logger,
	}
}

type watcherFactory struct {
	logger logger.Logger
}

func (factory *watcherFactory) Directory(folder string, checkIntervalInSeconds int, callback func(change *fswatch.FolderChange)) fswatch.Watcher {
	recurse := false

	var skipFunc = func(path string) bool {
		// don't skip
		return false
	}

	return factory.watchFolder(folder, checkIntervalInSeconds, recurse, skipFunc, callback)
}

func (factory *watcherFactory) SubDirectories(folder string, checkIntervalInSeconds int, callback func(change *fswatch.FolderChange)) fswatch.Watcher {
	recurse := false

	var skipFunc = func(path string) bool {
		// skip all files
		if isDirectory, _ := fsutil.IsDirectory(path); !isDirectory {
			return true
		}

		// skip all reserved directories
		return isReservedDirectory(path)
	}

	return factory.watchFolder(folder, checkIntervalInSeconds, recurse, skipFunc, callback)
}

func (factory *watcherFactory) AllFiles(folder string, checkIntervalInSeconds int, callback func(change *fswatch.FolderChange)) fswatch.Watcher {
	recurse := true

	var skipFunc = func(path string) bool {
		// don't skip anything
		return false
	}

	return factory.watchFolder(folder, checkIntervalInSeconds, recurse, skipFunc, callback)
}

func (factory *watcherFactory) Stop(folder string) {
	watcher, exists := watchers[folder]
	if !exists {
		return
	}

	watcher.Stop()
}

func (factory *watcherFactory) watchFolder(folder string, checkIntervalInSeconds int, recurse bool, skipFunc func(path string) bool, callback func(change *fswatch.FolderChange)) fswatch.Watcher {

	if existingWatcher, isReserved := factory.isReserved(folder); isReserved {
		factory.logger.Debug("Watcher %s already exists\n", folder)
		return existingWatcher
	}

	// look for changes in the item directory
	watcher := fswatch.NewFolderWatcher(folder, recurse, skipFunc, checkIntervalInSeconds)

	// start the watcher
	watcher.Start()

	defer func() {

		for watcher.IsRunning() {

			select {
			case change := <-watcher.ChangeDetails():
				callback(change)
			}

		}

		factory.release(folder)
		factory.logger.Debug("Exiting directory listener for folder %q.\n", folder)
	}()

	factory.reserve(folder, watcher)

	return watcher
}

func (factory *watcherFactory) reserve(folder string, watcher fswatch.Watcher) {
	watchers[folder] = watcher
}

func (factory *watcherFactory) release(folder string) {
	delete(watchers, folder)
}

func (factory *watcherFactory) isReserved(folder string) (fswatch.Watcher, bool) {
	if watcher, exists := watchers[folder]; exists {
		return watcher, true
	}

	return nil, false
}
