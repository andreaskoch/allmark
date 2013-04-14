// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"github.com/howeyc/fsnotify"
)

type FolderWatcher struct {
	Event chan *WatchEvent

	path              string
	watching          bool
	stopped           bool
	subFolderWatchers []*FolderWatcher
}

func NewFolderWatcher(folder string) (*FolderWatcher, error) {

	if ok, err := util.IsDirectory(folder); !ok {
		return nil, err // only handle files
	}

	return (&FolderWatcher{
		Event: make(chan *WatchEvent, 1),
		path:  folder,
	}).start(), nil
}

func (watcher *FolderWatcher) Stop() *FolderWatcher {
	watcher.watching = false
	watcher.stopped = true

	return watcher
}

func (watcher *FolderWatcher) Pause() *FolderWatcher {
	watcher.watching = false
	return watcher
}

func (watcher *FolderWatcher) IsWatching() bool {
	return watcher.watching
}

func (watcher *FolderWatcher) Resume() *FolderWatcher {
	watcher.watching = true
	return watcher
}

func (watcher *FolderWatcher) start() *FolderWatcher {

	// create watcher for all sub directories
	subDirectories := util.GetSubDirectories(watcher.path)
	watcher.subFolderWatchers = make([]*FolderWatcher, len(subDirectories), len(subDirectories))
	for index, subDirectory := range subDirectories {

		subFolderWatch, err := NewFolderWatcher(subDirectory)
		if err != nil {
			fmt.Errorf("Cannot create watch for folder %q.\nError: %s\n", subDirectory, err)
			continue
		}

		// dispatch event to the parent folder watcher
		go func() {
			for watcher.stopped == false {
				select {
				case event := <-subFolderWatch.Event:
					if watcher.IsWatching() {
						watcher.Event <- event
					}

				}
			}
		}()

		watcher.subFolderWatchers[index] = subFolderWatch
	}

	// start watching this folder
	watcher.watching = true

	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher.Pause()
	}

	go func() {
		for watcher.stopped == false {
			select {
			case event := <-fswatcher.Event:
				if watcher.IsWatching() {
					watcher.Event <- NewWatchEventFromFileEvent(event)
				}

			}
		}
	}()

	err = fswatcher.Watch(watcher.path)
	if err != nil {
		return watcher.Stop()
	}

	return watcher
}
