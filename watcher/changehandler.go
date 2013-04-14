// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type ChangeHandler struct {
	*Watcher

	callbacks CallbackList
}

func NewChangeHandler(path string) (*ChangeHandler, error) {

	// create a watcher
	folderwatcher, err := newWatcher(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to create folder watcher for %q.\nError: %s\n", path, err)
	}

	// create a new file change handler
	ChangeHandler := &ChangeHandler{
		Watcher:   folderwatcher,
		callbacks: NewCallbackList(),
	}

	// start watching
	ChangeHandler.startWatching()

	return ChangeHandler, err
}

func (changeHandler *ChangeHandler) startWatching() {
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

func (changeHandler *ChangeHandler) Throw(event *WatchEvent) {
	for _, callback := range changeHandler.callbacks.Values() {

		changeHandler.Pause()
		callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *ChangeHandler) OnChange(name string, callback ChangeHandlerCallback) {
	changeHandler.callbacks.Add(name, callback)
}
