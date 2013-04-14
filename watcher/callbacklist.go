// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

type ChangeHandlerCallback func(event *WatchEvent)

type CallbackList struct {
	names     map[string]int
	callbacks []ChangeHandlerCallback
}

func NewCallbackList() CallbackList {
	return CallbackList{
		names:     make(map[string]int),
		callbacks: make([]ChangeHandlerCallback, 0),
	}
}

func (list *CallbackList) Add(name string, callback ChangeHandlerCallback) {
	if position, ok := list.names[name]; ok {
		list.callbacks[position] = callback
		return
	}

	list.callbacks = append(list.callbacks, callback)
	positionOfNewCallback := len(list.callbacks) - 1
	list.names[name] = positionOfNewCallback
}

func (list *CallbackList) Values() []ChangeHandlerCallback {
	return list.callbacks
}
