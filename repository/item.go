// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package repository

import (
	"andyk/docs/filesystem"
	"os"
)

type Item struct {
	Path       string
	Files      []File
	ChildItems []Item
}

// Create a new repository item
func NewItem(path string, files []File, childItems []Item) Item {
	return Item{
		Path:       path,
		Files:      files,
		ChildItems: childItems,
	}
}

// Get all lines of a repository item
func (item *Item) GetLines() []string {
	inFile, err := os.Open(item.Path)
	if err != nil {
		panic("Could not read file.")
	}

	defer inFile.Close()

	return filesystem.GetLines(inFile)
}

// Render this repository item
func (item *Item) Render() {

	// render child items
	for _, child := range item.ChildItems {
		child.Render()
	}

	render := GetRenderer(item)
	render()
}

// Get a string representation of the current repository item
func (item Item) String() string {
	s := item.Path + "(" + item.GetHash() + ")\n"

	s += "\n"
	s += "Files:\n"
	if len(item.Files) > 0 {
		for _, file := range item.Files {
			s += " - " + file.Path + "\n"
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
