// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package indexer

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
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

func (item *Item) GetHash() string {
	itemBytes, readFileErr := ioutil.ReadFile(item.Path)
	if readFileErr != nil {
		return ""
	}

	sha1 := sha1.New()
	sha1.Write(itemBytes)

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:6]))
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
