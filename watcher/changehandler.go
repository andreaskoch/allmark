// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

type ChangeHandlerCallback func(event *WatchEvent)

type ChangeHandler interface {
	Throw(event *WatchEvent)

	OnCreate(name string, callback ChangeHandlerCallback)
	OnDelete(name string, callback ChangeHandlerCallback)
	OnModify(name string, callback ChangeHandlerCallback)
	OnRename(name string, callback ChangeHandlerCallback)
}

type CallbackKey struct {
	EventType    EventType
	CallbackName string
}

func NewCallbackKey(eventType EventType, callbackName string) CallbackKey {
	return CallbackKey{
		EventType:    eventType,
		CallbackName: callbackName,
	}
}

type CallbackEntry struct {
	EventType EventType
	Name      string
	Callback  ChangeHandlerCallback
}

func NewCallbackEntry(eventType EventType, name string, callback ChangeHandlerCallback) *CallbackEntry {
	return &CallbackEntry{
		EventType: eventType,
		Name:      name,
		Callback:  callback,
	}
}
