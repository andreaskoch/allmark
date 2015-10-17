// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Sitemap struct {
	Model
	Tree string
}

type SitemapEntry struct {
	Path        string         `json:"path"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Children      []SitemapEntry `json:"children"`
}
