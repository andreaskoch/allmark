// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type ToplevelNavigation struct {
	Entries []ToplevelEntry `json:"entries"`
}

func (navigation *ToplevelNavigation) IsAvailable() bool {
	return len(navigation.Entries) > 0
}

type ToplevelEntry struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
