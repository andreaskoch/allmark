// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FileChangeHandler struct {
	*FileWatcher

	callbacks map[string]ChangeHandlerCallback
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
	}

	// start watching
	fileChangeHandler.startWatching()

	return fileChangeHandler, err
}

func (changeHandler *FileChangeHandler) startWatching() {
	if changeHandler.callbacks == nil {
		changeHandler.callbacks = make(map[string]ChangeHandlerCallback) // initialize on first use
	}

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
	for _, callback := range changeHandler.callbacks {

		changeHandler.Pause()
		callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FileChangeHandler) OnChange(name string, callback ChangeHandlerCallback) {
	if _, ok := changeHandler.callbacks[name]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[name] = callback
}
