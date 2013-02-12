// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package model

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type RepositoryItem struct {
	Path       string
	Files      []RepositoryItemFile
	ChildItems []RepositoryItem
	Type       string
}

// Create a new repository item
func NewRepositoryItem(itemType string, path string, files []RepositoryItemFile, childItems []RepositoryItem) RepositoryItem {
	return RepositoryItem{
		Path:       path,
		Files:      files,
		ChildItems: childItems,
		Type:       itemType,
	}
}

// Render this repository item
func (item *RepositoryItem) Render() {

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

	doc := item.getParsedDocument()

	content := "<!-- " + item.GetHash() + " -->"
	content += "\nTitle: " + doc["title"]

	_ = ioutil.WriteFile(renderedItemPath, []byte(content), 0644)
}

func (item *RepositoryItem) getParsedDocument() map[string]string {
	bytes := item.getBytes()
	doc := make(map[string]string)

	doc["title"] = getTitle(bytes)
	doc["description"] = getDescription(bytes)

	return doc
}

func getTitle(content []byte) string {
	buffer := bytes.NewBuffer(content)

	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")
	for line, err := buffer.ReadString(10); err != io.EOF; {
		matches := titleRegexp.FindStringSubmatch(line)

		if len(matches) == 2 {
			return matches[1]
		}
	}

	return "No Title"
}

func getDescription(content []byte) string {
	buffer := bytes.NewBuffer(content)

	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")
	for line, err := buffer.ReadString(10); err != io.EOF; {
		matches := titleRegexp.FindStringSubmatch(line)

		if len(matches) == 2 {
			return matches[1]
		}
	}

	return "No Title"
}

// Get the bytes of this repository item
func (item *RepositoryItem) getBytes() []byte {
	bytes, err := ioutil.ReadFile(item.Path)
	if err != nil {
		return make([]byte, 0)
	}

	return bytes
}

// Get the hash code of the rendered item
func (item *RepositoryItem) GetRenderedItemHash() string {
	renderedItemPath := item.GetRenderedItemPath()

	file, err := os.Open(renderedItemPath)
	if err != nil {
		// file does not exist or cannot be accessed
		return ""
	}

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
func (item *RepositoryItem) GetRenderedItemPath() string {
	itemDirectory := filepath.Dir(item.Path)
	renderedFilePath := filepath.Join(itemDirectory, item.Type+".html")
	return renderedFilePath
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

// Get a string representation of the current repository item
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
