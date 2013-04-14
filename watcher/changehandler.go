// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

type ChangeHandler interface {
	Throw(event *WatchEvent)
	OnChange(name string, callback ChangeHandlerCallback)
}
