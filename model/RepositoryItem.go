// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package model

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type RepositoryItem struct {
	Path       string
	Files      []RepositoryItemFile
	ChildItems []RepositoryItem
	Type       string
}

func NewRepositoryItem(itemType string, path string, files []RepositoryItemFile, childItems []RepositoryItem) RepositoryItem {
	return RepositoryItem{
		Path:       path,
		Files:      files,
		ChildItems: childItems,
		Type:       itemType,
	}
}

func (item *RepositoryItem) Render() {

	// render child items
	for _, child := range item.ChildItems {
		child.Render()
	}

	// assemble file path of the rendered html file
	itemDirectory := filepath.Dir(item.Path)
	renderedFilePath := filepath.Join(itemDirectory, item.Type+".html")

	// create html file if it does not exist
	if _, getFileStatError := os.Stat(renderedFilePath); getFileStatError != nil {
		_ = ioutil.WriteFile(renderedFilePath, []byte(""), 0644)
	}
}

func (item *RepositoryItem) GetHash() string {
	itemBytes, readFileErr := ioutil.ReadFile(item.Path)
	if readFileErr != nil {
		return ""
	}

	sha1 := sha1.New()
	sha1.Write(itemBytes)

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:6]))
}

func (item *RepositoryItem) String() string {
	s := item.Path + "(Type: " + item.Type + ", Hash: " + item.GetHash() + ")\n"

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
