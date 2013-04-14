// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FolderChangeHandler struct {
	*FolderWatcher

	callbacks map[CallbackKey]*CallbackEntry
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
		changeHandler.callbacks = make(map[CallbackKey]*CallbackEntry) // initialize on first use
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
	fmt.Printf("%s: %s\n", event.Type, event.Filepath)
	for _, entry := range changeHandler.getHandlersByType(event.Type) {

		changeHandler.Pause()
		entry.Callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FolderChangeHandler) OnCreate(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(CREATED, name, callback)
}

func (changeHandler *FolderChangeHandler) OnDelete(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(DELETED, name, callback)
}

func (changeHandler *FolderChangeHandler) OnModify(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(MODIFIED, name, callback)
}

func (changeHandler *FolderChangeHandler) OnRename(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(RENAMED, name, callback)
}

func (changeHandler *FolderChangeHandler) addHandler(eventType EventType, name string, callback ChangeHandlerCallback) {

	key := NewCallbackKey(eventType, name)

	if _, ok := changeHandler.callbacks[key]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[key] = NewCallbackEntry(eventType, name, callback)
}

func (changeHandler *FolderChangeHandler) getHandlersByType(eventType EventType) []*CallbackEntry {

	entries := make([]*CallbackEntry, 0, len(changeHandler.callbacks))

	for key, entry := range changeHandler.callbacks {
		if key.EventType == eventType {
			entries = append(entries, entry)
		}
	}

	return entries
}
