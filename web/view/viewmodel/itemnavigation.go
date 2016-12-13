// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type ItemNavigation struct {
	Parent   NavEntry `json:"parent"`
	Previous NavEntry `json:"previous"`
	Next     NavEntry `json:"next"`
}

// IsAvailable returns a flag indicating whether the item navigation model is initialized or not.
func (nav ItemNavigation) IsAvailable() bool {
	return nav.Parent.Path != "" || nav.Previous.Path != "" || nav.Next.Path != ""
}

type NavEntry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Path        string `json:"path"`
}
