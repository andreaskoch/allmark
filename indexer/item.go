// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package indexer

import (
	"andyk/docs/util"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	UnknownItemType    = "unknown"
	DocumentItemType   = "document"
	MessageItemType    = "message"
	LocationItemType   = "location"
	CommentItemType    = "message"
	TagItemType        = "tag"
	RepositoryItemType = "repository"
)

type Item struct {
	Path         string
	RenderedPath string
	Files        []File
	ChildItems   []Item
	Type         string
}

// Create a new repository item
func NewItem(path string, files []File, childItems []Item) Item {

	var getItemType = func() string {
		filename := filepath.Base(path)
		return getItemTypeFromFilename(filename)
	}

	return Item{
		Path:         path,
		RenderedPath: getRenderedItemPath(path),
		Files:        files,
		ChildItems:   childItems,
		Type:         getItemType(),
	}
}

func (item Item) GetFilename() string {
	return filepath.Base(item.Path)
}

func (item Item) GetHash() string {
	itemBytes, readFileErr := ioutil.ReadFile(item.Path)
	if readFileErr != nil {
		return ""
	}

	sha1 := sha1.New()
	sha1.Write(itemBytes)

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:6]))
}

func (item Item) Walk(walkFunc func(item Item)) {

	walkFunc(item)

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func (item Item) IsRendered() bool {
	return util.FileExists(item.RenderedPath)
}

func (item Item) GetAbsolutePath() string {
	return item.RenderedPath
}

func (item Item) GetRelativePath(basePath string) string {

	fullItemPath := item.RenderedPath
	relativePath := strings.Replace(fullItemPath, basePath, "", 1)
	relativePath = "/" + strings.TrimLeft(relativePath, "/")

	return relativePath
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(itemPath string) string {
	itemDirectory := filepath.Dir(itemPath)
	renderedFilePath := filepath.Join(itemDirectory, "index.html")
	return renderedFilePath
}

func getItemTypeFromFilename(filename string) string {

	lowercaseFilename := strings.ToLower(filename)

	switch lowercaseFilename {
	case "document.md":
		return DocumentItemType

	case "message.md":
		return MessageItemType

	case "location.md":
		return LocationItemType

	case "comment.md":
		return CommentItemType

	case "tag.md":
		return TagItemType

	case "repository.md":
		return RepositoryItemType
	}

	return UnknownItemType
}
