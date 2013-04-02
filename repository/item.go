// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package model defines the basic
	data structures of the docs engine.
*/
package repository

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/util"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	UnknownItemType      = "unknown"
	DocumentItemType     = "document"
	PresentationItemType = "presentation"
	CollectionItemType   = "collection"
	MessageItemType      = "message"
	ImageGalleryItemType = "imagegallery"
	LocationItemType     = "location"
	CommentItemType      = "comment"
	TagItemType          = "tag"
	RepositoryItemType   = "repository"
)

type Item struct {
	Title        string
	Description  string
	Content      string
	Path         string
	RenderedPath string
	Files        []*File
	ChildItems   []*Item
	MetaData     MetaData
	Type         string

	onChangeCallbacks  map[string]func(item *Item)
	itemIsBeingWatched bool
}

// Create a new repository item
func NewItem(path string, childItems []*Item) (item *Item, err error) {

	itemType := getItemType(path)

	if itemType == UnknownItemType {
		err = errors.New(fmt.Sprintf("The item %q does not match any of the known item types.", path))
	}

	item = &Item{
		Path:         path,
		RenderedPath: getRenderedItemPath(path),
		ChildItems:   childItems,
		Type:         itemType,
	}

	item.IndexFiles()
	item.startWatch()

	return item, err
}

func (item *Item) String() string {
	return fmt.Sprintf("Item %s\n", item.Path)
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

func (item *Item) Walk(walkFunc func(item *Item)) {

	walkFunc(item)

	// add all children
	for _, child := range item.ChildItems {
		child.Walk(walkFunc)
	}
}

func (item Item) GetFolder() string {
	return filepath.Dir(item.Path)
}

func (item Item) IsRendered() bool {
	return util.FileExists(item.RenderedPath)
}

func (item *Item) Directory() string {
	return filepath.Dir(item.Path)
}

func (item Item) AbsolutePath() string {
	return item.RenderedPath
}

func (item Item) RelativePath(basePath string) string {

	pathSeperator := string(os.PathSeparator)
	fullItemPath := item.RenderedPath
	relativePath := strings.Replace(fullItemPath, basePath, "", 1)
	relativePath = pathSeperator + strings.TrimLeft(relativePath, pathSeperator)

	return relativePath
}

func (item *Item) Render(renderFunc func(item *Item) *Item) {
	item.pauseWatch()
	defer item.resumeWatch()

	renderFunc(item)
}

func (item *Item) RegisterOnChangeCallback(name string, callbackFunction func(item *Item)) {

	if item.onChangeCallbacks == nil {
		item.onChangeCallbacks = make(map[string]func(item *Item))
	}

	if _, ok := item.onChangeCallbacks[name]; ok {
		fmt.Printf("Change callback %q already present.", name)
	}

	item.onChangeCallbacks[name] = callbackFunction
}

func (item *Item) IndexFiles() *Item {

	filesDirectory := filepath.Join(item.Directory(), "files")
	itemFiles := make([]*File, 0, 5)
	filesDirectoryEntries, _ := ioutil.ReadDir(filesDirectory)

	for _, file := range filesDirectoryEntries {
		if file.IsDir() {
			continue
		}

		absoluteFilePath := filepath.Join(filesDirectory, file.Name())
		repositoryFile := NewFile(absoluteFilePath)

		itemFiles = append(itemFiles, repositoryFile)
	}

	item.Files = itemFiles
	return item
}

func (item *Item) pauseWatch() {
	fmt.Printf("Pausing watch on item %s\n", item)
	item.itemIsBeingWatched = true
}

func (item *Item) watchIsPaused() bool {
	return item.itemIsBeingWatched
}

func (item *Item) resumeWatch() {
	fmt.Printf("Resuming watch on item %s\n", item)
	item.itemIsBeingWatched = false
}

func (item *Item) startWatch() *Item {

	item.itemIsBeingWatched = true

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error while creating watch for item %q. Error: %v", item, err)
		return item
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:

				if !item.watchIsPaused() {
					fmt.Println("Item changed ->", event)

					for name, callback := range item.onChangeCallbacks {
						fmt.Printf("Item changed. Executing callback %q on for item %q\n", name, item)
						callback(item)
					}
				}

			case err := <-watcher.Error:
				fmt.Printf("Watch error on item %q. Error: %v\n", item, err)
			}
		}
	}()

	err = watcher.Watch(item.Path)
	if err != nil {
		fmt.Printf("Error while creating watch for folder %q. Error: %v\n", item.Path, err)
	}

	return item
}

// Get the item type from the given item path
func getItemType(itemPath string) string {
	filename := filepath.Base(itemPath)
	return getItemTypeFromFilename(filename)
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
	case "document.md", "readme.md":
		return DocumentItemType

	case "presentation.md":
		return PresentationItemType

	case "collection.md":
		return CollectionItemType

	case "message.md":
		return MessageItemType

	case "imagegallery.md":
		return ImageGalleryItemType

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
