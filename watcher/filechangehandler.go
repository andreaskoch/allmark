// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FileChangeHandler struct {
	*FileWatcher

	callbacks CallbackList
}

func NewFileChangeHandler(filePath string) (*FileChangeHandler, error) {

	// create a watcher
	filewatcher, err := NewFileWatcher(filePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to create file watcher for %q.\nError: %s\n", filePath, err)
	}

	// create a new file change handler
	fileChangeHandler := &FileChangeHandler{
		FileWatcher: filewatcher,
		callbacks:   NewCallbackList(),
	}

	// start watching
	fileChangeHandler.startWatching()

	return fileChangeHandler, err
}

func (changeHandler *FileChangeHandler) startWatching() {
	// start watching for changes
	go func() {
		for {
			select {
			case event := <-changeHandler.Event:
				changeHandler.Throw(event)
			}
		}
	}()
}

func (changeHandler *FileChangeHandler) Throw(event *WatchEvent) {
	for _, callback := range changeHandler.callbacks.Values() {

		changeHandler.Pause()
		callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FileChangeHandler) OnChange(name string, callback ChangeHandlerCallback) {
	changeHandler.callbacks.Add(name, callback)
}
