// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FolderChangeHandler struct {
	*FolderWatcher

	callbacks CallbackList
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
		callbacks:     NewCallbackList(),
	}

	// start watching
	FolderChangeHandler.startWatching()

	return FolderChangeHandler, err
}

func (changeHandler *FolderChangeHandler) startWatching() {
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
	for _, callback := range changeHandler.callbacks.Values() {

		changeHandler.Pause()
		callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FolderChangeHandler) OnChange(name string, callback ChangeHandlerCallback) {
	changeHandler.callbacks.Add(name, callback)
}
