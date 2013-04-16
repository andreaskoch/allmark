// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"github.com/howeyc/fsnotify"
)

type Watcher struct {
	Event chan *WatchEvent

	path        string
	watching    bool
	stopped     bool
	subWatchers []*Watcher
}

func newWatcher(path string) (*Watcher, error) {
	return (&Watcher{
		Event: make(chan *WatchEvent, 1),
		path:  path,
	}).start(), nil
}

func (watcher *Watcher) Stop() *Watcher {
	watcher.watching = false
	watcher.stopped = true

	return watcher
}

func (watcher *Watcher) Pause() *Watcher {
	watcher.watching = false
	return watcher
}

func (watcher *Watcher) IsWatching() bool {
	return watcher.watching
}

func (watcher *Watcher) Resume() *Watcher {
	watcher.watching = true
	return watcher
}

func (watcher *Watcher) start() *Watcher {

	// create watcher for all sub directories
	if ok, _ := util.IsDirectory(watcher.path); ok {

		subDirectories := util.GetSubDirectories(watcher.path)
		watcher.subWatchers = make([]*Watcher, 0, len(subDirectories))
		for _, subDirectory := range subDirectories {

			subFolderWatch, err := newWatcher(subDirectory)
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

			watcher.subWatchers = append(watcher.subWatchers, subFolderWatch)
		}

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
					watcher.Event <- newWatchEventFromFileEvent(event)
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
