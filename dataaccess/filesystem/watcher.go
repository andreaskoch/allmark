// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/common/route"
	"fmt"
	"github.com/andreaskoch/go-fswatch"
)

type watcherPather interface {
	Path() string
	IsDirectory() bool
	Recurse() bool
}

// A file path
type watcherFilePath struct {
	path string
}

func (w watcherFilePath) Path() string {
	return w.path
}

func (w watcherFilePath) IsDirectory() bool {
	return false
}

func (w watcherFilePath) Recurse() bool {
	return false
}

// A directory path
type watcherDirectoryPath struct {
	path    string
	recurse bool
}

func (w watcherDirectoryPath) Path() string {
	return w.path
}

func (w watcherDirectoryPath) IsDirectory() bool {
	return true
}

func (w watcherDirectoryPath) Recurse() bool {
	return w.recurse
}

func newFilesystemWatcher(logger logger.Logger) *filesystemWatcher {
	return &filesystemWatcher{
		logger:   logger,
		watchers: make(map[string][]fswatch.Watcher),
	}
}

type filesystemWatcher struct {
	logger   logger.Logger
	watchers map[string][]fswatch.Watcher
}

func (watcher *filesystemWatcher) Start(route route.Route, watcherPaths []watcherPather) (chan bool, error) {

	// check if there are already watchers
	if _, exists := watcher.watchers[routeToString(route)]; exists {
		return nil, fmt.Errorf("The watchers for route %q are already running.", route.String())
	}

	watcher.logger.Debug("Starting to watch %q", route.String())

	backChannel := make(chan bool, 1)

	// create a watcher for every path
	watchers := make([]fswatch.Watcher, 0)
	for _, watcherPath := range watcherPaths {

		// create a filesystem watcher
		var pathWatcher fswatch.Watcher
		if watcherPath.IsDirectory() {
			pathWatcher = watcher.createDirectoryWatcher(watcherPath.Path(), watcherPath.Recurse(), backChannel)
		} else {
			pathWatcher = watcher.createFileWatcher(watcherPath.Path(), backChannel)
		}

		watchers = append(watchers, pathWatcher)
	}

	// store the list
	watcher.watchers[routeToString(route)] = watchers

	return backChannel, nil
}

func (watcher *filesystemWatcher) Stop(route route.Route) {

	// Get the requested watcher list
	watcherList, exists := watcher.watchers[routeToString(route)]
	if !exists {
		return
	}

	watcher.logger.Debug("Stopping to watch %q", route.String())

	// stop all watchers
	for _, listEntry := range watcherList {
		listEntry.Stop()
	}

	// remove from list
	delete(watcher.watchers, routeToString(route))
}

func (watcher *filesystemWatcher) IsRunning(route route.Route) bool {
	_, exists := watcher.watchers[routeToString(route)]
	return exists
}

func (watcher *filesystemWatcher) createFileWatcher(filePath string, backChannel chan bool) fswatch.Watcher {

	checkIntervalInSeconds := 1
	filewatcher := fswatch.NewFileWatcher(filePath, checkIntervalInSeconds)
	filewatcher.Start()

	// the go-routine which waits for changes
	go func() {
		running := true
		for running {

			select {
			case <-filewatcher.Modified():
				backChannel <- true

			case <-filewatcher.Moved():
				running = false

			case <-filewatcher.Stopped():
				running = false
			}
		}
	}()

	return filewatcher
}

func (watcher *filesystemWatcher) createDirectoryWatcher(directoryPath string, recurse bool, backChannel chan bool) fswatch.Watcher {
	checkIntervalInSeconds := 1

	skipNoFiles := func(path string) bool {
		return false
	}

	folderWatcher := fswatch.NewFolderWatcher(directoryPath, recurse, skipNoFiles, checkIntervalInSeconds)
	folderWatcher.Start()

	// the go-routine which waits for changes
	go func() {
		running := true
		for running {

			select {
			case <-folderWatcher.Modified():
				backChannel <- true

			case <-folderWatcher.Moved():
				running = false

			case <-folderWatcher.Stopped():
				running = false
			}
		}
	}()

	return folderWatcher
}

func routeToString(route route.Route) string {
	return route.Value()
}
