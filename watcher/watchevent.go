// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"github.com/howeyc/fsnotify"
)

type WatchEvent struct {
	Filepath string
}

func NewWatchEventFromFileEvent(event *fsnotify.FileEvent) *WatchEvent {
	return NewWatchEvent(event.Name)
}

func NewWatchEvent(filepath string) *WatchEvent {
	return &WatchEvent{
		Filepath: filepath,
	}
}

func (watchEvent *WatchEvent) String() string {
	return fmt.Sprintf("Event (%s)", watchEvent.Filepath)
}
