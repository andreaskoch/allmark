// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Search struct {
	Query        string
	Page         int
	ItemsPerPage int
	Results      []SearchResult
}

type SearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
	PubDate     string `json:"pubDate"`
}
