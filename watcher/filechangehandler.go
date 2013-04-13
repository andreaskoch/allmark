// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"strings"
)

type CallbackEntry struct {
	EventType string
	Name      string
	Callback  ChangeHandlerCallback
}

func NewCallbackEntry(eventType, name string, callback ChangeHandlerCallback) *CallbackEntry {

	return &CallbackEntry{
		EventType: eventType,
		Name:      name,
		Callback:  callback,
	}

}

type FileChangeHandler struct {
	*FileWatcher

	callbacks map[string]*CallbackEntry
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

func (changeHandler *FileChangeHandler) OnCreate(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Create", name, callback)
}

func (changeHandler *FileChangeHandler) OnDelete(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Delete", name, callback)
}

func (changeHandler *FileChangeHandler) OnModify(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Modify", name, callback)
}

func (changeHandler *FileChangeHandler) OnRename(name string, callback ChangeHandlerCallback) {
	changeHandler.addHandler("Rename", name, callback)
}

func (changeHandler *FileChangeHandler) addHandler(eventType, name string, callback ChangeHandlerCallback) {

	key := fmt.Sprintf("%s - %s", eventType, name)

	if _, ok := changeHandler.callbacks[key]; ok {
		fmt.Printf("WARNING: Change callback %q already present.", name)
	}

	changeHandler.callbacks[key] = NewCallbackEntry(eventType, name, callback)
}
