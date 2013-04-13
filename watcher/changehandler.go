// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

type ChangeHandlerCallback func(event *WatchEvent)

type ChangeHandler interface {
	OnCreate(name string, callback ChangeHandlerCallback)
	OnDelete(name string, callback ChangeHandlerCallback)
	OnModify(name string, callback ChangeHandlerCallback)
	OnRename(name string, callback ChangeHandlerCallback)
}
