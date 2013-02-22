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
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Item struct {
	Path       string
	Files      []File
	ChildItems []Item
	Type       string
}

// Create a new repository item
func NewItem(itemType string, path string, files []File, childItems []Item) Item {
	return Item{
		Path:       path,
		Files:      files,
		ChildItems: childItems,
		Type:       itemType,
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

	// the path of the rendered repostory item
	renderedItemPath := item.GetRenderedItemPath()

	// check if rendering is required
	itemHashCode := item.GetHash()
	renderedItemHashCode := item.GetRenderedItemHash()

	// Abort if the hash has not changed
	if itemHashCode == renderedItemHashCode {
		return
	}

	renderer := GetRenderer(item)
	renderer.Execute()
}

// Get the hash code of the rendered item
func (item Item) GetRenderedItemHash() string {
	renderedItemPath := item.GetRenderedItemPath()

	file, err := os.Open(renderedItemPath)
	if err != nil {
		// file does not exist or cannot be accessed
		return ""
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	firstLineBytes, _ := fileReader.ReadBytes('\n')
	if firstLineBytes == nil {
		// first line cannot be read
		return ""
	}

	// extract hash from line
	hashCodeRegexp := regexp.MustCompile("<!-- (\\w+) -->")
	matches := hashCodeRegexp.FindStringSubmatch(string(firstLineBytes))
	if len(matches) != 2 {
		return ""
	}

	extractedHashcode := matches[1]

	return string(extractedHashcode)
}

// Get the filepath of the rendered repository item
func (item Item) GetRenderedItemPath() string {
	itemDirectory := filepath.Dir(item.Path)
	renderedFilePath := filepath.Join(itemDirectory, item.Type+".html")
	return renderedFilePath
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
