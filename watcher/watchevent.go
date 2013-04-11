// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import "fmt"

type WatchEvent struct {
	Filepath string
	Type     EventType
}

func NewWatchEvent(filepath, eventType string) *WatchEvent {
	return &WatchEvent{
		Filepath: filepath,
		Type:     EventTypeFromText(eventType),
	}
}

func (watchEvent *WatchEvent) String() string {
	return fmt.Sprintf("Event (Type: %s, Path: %s)", watchEvent.Type, watchEvent.Filepath)
}
