// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"strings"
)

type FolderChangeHandler struct {
	*FolderWatcher

	callbacks map[string]*CallbackEntry
}

func NewFolderChangeHandler(path string) (*FolderChangeHandler, error) {

	// create a watcher
	folderwatcher, err := NewFolderWatcher(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to create file watcher for %q.\nError: %s\n", path, err)
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
		changeHandler.callbacks = make(map[string]*CallbackEntry) // initialize on first use
	}

	// start watching for changes
	go func() {
		for {
			select {
			case event := <-changeHandler.Event:

				fmt.Printf("%s: %s\n", strings.ToUpper(event.Type.String()), event.Filepath)
				for _, entry := range changeHandler.callbacks {

					changeHandler.Pause()
					entry.Callback(event)
					changeHandler.Resume()

				}
			}
		}
	}()
}

func (changeHandler *FolderChangeHandler) OnCreate(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Create", name, callback)
}

func (changeHandler *FolderChangeHandler) OnDelete(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Delete", name, callback)
}

func (changeHandler *FolderChangeHandler) OnModify(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Modify", name, callback)
}

func (changeHandler *FolderChangeHandler) OnRename(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Rename", name, callback)
}

func (changeHandler *FolderChangeHandler) addHandler(eventType, name string, callback ChangeHandlerCallback) {

	key := fmt.Sprintf("%s - %s", eventType, name)

	if _, ok := changeHandler.callbacks[key]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[key] = NewCallbackEntry(eventType, name, callback)
}
