// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"github.com/andreaskoch/allmark/util"
	"github.com/howeyc/fsnotify"
)

type FileWatcher struct {
	Event chan *WatchEvent

	path     string
	watching bool
	stopped  bool
}

func NewFileWatcher(filepath string) (*FileWatcher, error) {

	if ok, err := util.IsFile(filepath); !ok {
		return nil, err // only handle files
	}

	return (&FileWatcher{
		Event: make(chan *WatchEvent, 1),
		path:  filepath,
	}).start(), nil
}

func (watcher *FileWatcher) Stop() *FileWatcher {
	watcher.watching = false
	watcher.stopped = true

	return watcher
}

func (watcher *FileWatcher) Pause() *FileWatcher {
	watcher.watching = false
	return watcher
}

func (watcher *FileWatcher) IsWatching() bool {
	return watcher.watching
}

func (watcher *FileWatcher) Resume() *FileWatcher {
	watcher.watching = true
	return watcher
}

func (watcher *FileWatcher) start() *FileWatcher {

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
					watcher.Event <- getWatchEventFromFileEvent(event)
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
