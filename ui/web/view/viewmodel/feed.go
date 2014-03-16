// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type FeedEntry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"route"`
	PubDate     string `json:"pubDate"`
}
