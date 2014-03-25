// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Search struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`

	Page         int `json:"page"`
	ItemsPerPage int `json:"itemPerPage"`

	ResultCount      int `json:"resultCount"`
	TotalResultCount int `json:"totalResultCount"`
}

type SearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
	Path        string `json:"path"`
}
