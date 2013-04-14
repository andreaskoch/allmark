// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
)

type FileChangeHandler struct {
	*FileWatcher

	callbacks map[CallbackKey]*CallbackEntry
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

func (changeHandler *FileChangeHandler) Throw(event *WatchEvent) {
	fmt.Printf("%s: %s\n", event.Type, event.Filepath)
	for _, entry := range changeHandler.getHandlersByType(event.Type) {

		changeHandler.Pause()
		entry.Callback(event)
		changeHandler.Resume()

	}
}

func (changeHandler *FileChangeHandler) OnCreate(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(CREATED, name, callback)
}

func (changeHandler *FileChangeHandler) OnDelete(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(DELETED, name, callback)
}

func (changeHandler *FileChangeHandler) OnModify(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(MODIFIED, name, callback)
}

func (changeHandler *FileChangeHandler) OnRename(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler(RENAMED, name, callback)
}

func (changeHandler *FileChangeHandler) addHandler(eventType EventType, name string, callback ChangeHandlerCallback) {

	key := NewCallbackKey(eventType, name)

	if _, ok := changeHandler.callbacks[key]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[key] = NewCallbackEntry(eventType, name, callback)
}

func (changeHandler *FileChangeHandler) getHandlersByType(eventType EventType) []*CallbackEntry {

	entries := make([]*CallbackEntry, 0, len(changeHandler.callbacks))

	for key, entry := range changeHandler.callbacks {
		if key.EventType == eventType {
			entries = append(entries, entry)
		}
	}

	return entries
}
