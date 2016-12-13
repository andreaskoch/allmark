// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type TypeAhead struct {
	Index int `json:"index"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
	Path        string `json:"path"`

	Value  string   `json:"value"`
	Tokens []string `json:"tokens"`
}
