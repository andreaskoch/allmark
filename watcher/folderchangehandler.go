// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FolderChangeHandler struct {
	*FolderWatcher

	callbacks map[string]ChangeHandlerCallback
}

func NewFolderChangeHandler(path string) (*FolderChangeHandler, error) {

	// create a watcher
	folderwatcher, err := NewFolderWatcher(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to create folder watcher for %q.\nError: %s\n", path, err)
	}

	// create a new file change handler
	FolderChangeHandler := &FolderChangeHandler{
		FolderWatcher: folderwatcher,
	}

	// start watching
	FolderChangeHandler.startWatching()

	return FolderChangeHandler, err
}

func (changeHandler *FolderChangeHandler) startWatching() {
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

func (changeHandler *FolderChangeHandler) Throw(event *WatchEvent) {
	for _, callback := range changeHandler.callbacks {

		changeHandler.Pause()
		callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FolderChangeHandler) OnChange(name string, callback ChangeHandlerCallback) {
	if _, ok := changeHandler.callbacks[name]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[name] = callback
}
