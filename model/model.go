// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package model

import "time"

type RepositoryItem struct {
	Path string
}

func NewRepositoryItem(path string) *RepositoryItem {
	return &RepositoryItem{Path: path}
}

type Document struct {
	Path    string // The documents folder
	Content string // The document content

	Title       string // The document title
	Description string // A short description of the document content.

	// Meta information
	Language string    // [optional] The ISO language code document (e.g. "en-GB", "de-DE")
	Date     time.Time // [optional] The date the document has been created
}
