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
	Path       string
	Files      []string
	ChildItems []RepositoryItem
}

func NewRepositoryItem(path string, files []string, childItems []RepositoryItem) RepositoryItem {
	return RepositoryItem{
		Path:       path,
		Files:      files,
		ChildItems: childItems,
	}
}

func (item *RepositoryItem) String() string {
	s := item.Path + "\n"

	s += "\n"
	s += "Files:\n"
	if len(item.Files) > 0 {
		for _, file := range item.Files {
			s += " - " + file + "\n"
		}
	} else {
		s += "<none>\n"
	}

	s += "\n"
	s += "ChildItems:\n"
	if len(item.ChildItems) > 0 {
		for _, child := range item.ChildItems {
			s += child.String()
		}
	} else {
		s += "<none>\n"
	}
	s += "\n"

	return s
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
