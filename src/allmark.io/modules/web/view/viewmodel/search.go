// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Search struct {
	Model
	Results SearchResults
}

type SearchResults struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`

	Page         int `json:"page"`
	ItemsPerPage int `json:"itemPerPage"`

	StartIndex       int `json:"startIndex"`
	ResultCount      int `json:"resultCount"`
	TotalResultCount int `json:"totalResultCount"`
}

type SearchResult struct {
	Index int `json:"index"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
	Path        string `json:"path"`
}
