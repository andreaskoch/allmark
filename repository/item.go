// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"github.com/andreaskoch/allmark/view"
	"github.com/andreaskoch/allmark/watcher"
	"io/ioutil"
	"path/filepath"
)

const (
	FilesDirectoryName = "files"
)

type Item struct {
	*view.Model

	Modified chan bool
	Moved    chan bool

	Level  int
	Files  *FileIndex
	Childs []*Item

	directory string
	path      string
	isVirtual bool

	pathProvider *path.Provider

	changeFuncs []func(item *Item)
}

func NewRoot(indexPath string) (item *Item, err error) {

	// check if path exists
	if !util.PathExists(indexPath) {
		return nil, fmt.Errorf("The path %q does not exist.", indexPath)
	}

	if isReservedDirectory(indexPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be a root.", indexPath)
	}

	// check if the path is a directory
	if isDirectory, _ := util.IsDirectory(indexPath); !isDirectory {
		indexPath = filepath.Dir(indexPath)
	}

	rootItem, err := newItem(indexPath, 0)
	if err != nil {
		return nil, err
	}

	return rootItem, nil
}

func newItem(itemPath string, level int) (item *Item, err error) {

	// abort if path does not exist
	if !util.PathExists(itemPath) {
		return nil, fmt.Errorf("The path %q does not exist.", itemPath)
	}

	// abort if path is reserved
	if isReservedDirectory(itemPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemPath)
	}

	// check if its a virtual item or a markdown item
	isVirtualItem := false
	itemDirectory := filepath.Dir(itemPath)

	if isDirectory, _ := util.IsDirectory(itemPath); isDirectory {

		if found, filepath := findMarkdownFileInDirectory(itemPath); found {

			itemDirectory = itemPath
			itemPath = filepath

		} else {

			isVirtualItem = true
			itemDirectory = itemPath

		}

	} else if !markdown.IsMarkdownFile(itemPath) {

		// the supplied item path does not point to a markdown file
		return nil, fmt.Errorf("%q is not a markdown file.", itemPath)
	}

	// create a path provider
	pathProviderDirectory := itemDirectory
	if level > 0 {
		pathProviderDirectory = filepath.Dir(itemDirectory)
	}

	pathProvider := path.NewProvider(pathProviderDirectory, false)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, FilesDirectoryName)
	fileIndex, err := NewFileIndex(filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item = &Item{
		Modified: make(chan bool),
		Moved:    make(chan bool),

		Level: level,
		Files: fileIndex,

		pathProvider: pathProvider,
		directory:    itemDirectory,
		path:         itemPath,
		isVirtual:    isVirtualItem,
		changeFuncs:  make([]func(item *Item), 0),
	}

	// find childs
	item.updateChilds()

	// look for item changes
	if !isVirtualItem {
		go func() {
			fileWatcher := watcher.NewFileWatcher(itemPath).Start()

			for fileWatcher.IsRunning() {

				select {
				case <-fileWatcher.Modified:

					item.Modified <- true

				case <-fileWatcher.Moved:

					item.Moved <- true
				}

			}
		}()
	}

	// look for file changes
	go func() {
		for {
			select {
			case <-item.Files.Changed:
				fmt.Println("Item changed because files changed")
				item.Modified <- true

			}
		}
	}()

	return item, nil
}

func (item *Item) OnChange(name string, expr func(i *Item)) {
	item.changeFuncs = append(item.changeFuncs, expr)
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.path)
}

func (item *Item) Path() string {
	return item.path
}

func (item *Item) Directory() string {
	return item.directory
}

func (item *Item) PathType() string {
	return path.PatherTypeItem
}

func (item *Item) IsVirtual() bool {
	return item.isVirtual
}

func (item *Item) PathProvider() *path.Provider {
	return item.pathProvider
}

func (item *Item) HasChild(itemDirectory string) bool {
	if item.Childs == nil {
		return false
	}

	for _, childItem := range item.Childs {
		if childItem.Directory() == itemDirectory {
			return true
		}
	}

	return false
}

func (item *Item) updateChilds() {

	childItemDirectories := item.getChildItemDirectories()

	if item.Childs == nil {
		item.Childs = make([]*Item, 0, len(childItemDirectories))
	}

	// add new childs
	for _, childItemDirectory := range childItemDirectories {

		if item.HasChild(childItemDirectory) {
			continue // Child already present
		}

		if child, err := newItem(childItemDirectory, item.Level+1); err == nil {
			item.Childs = append(item.Childs, child) // append new child
		} else {
			fmt.Printf("Could not create a item for folder %q. Error: %s\n", childItemDirectory, err)
		}
	}

	// remove deleted childs
	newChildList := make([]*Item, 0, len(childItemDirectories))
	for _, child := range item.Childs {
		if util.SliceContainsElement(childItemDirectories, child.Directory()) {
			newChildList = append(newChildList, child)
		}
	}

	item.Childs = newChildList
}

func (item *Item) getChildItemDirectories() []string {

	directory := item.Directory()
	directories := make([]string, 0)

	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		if !entry.IsDir() {
			continue // skip files
		}

		childDirectory := filepath.Join(directory, entry.Name())
		if isReservedDirectory(childDirectory) {
			continue // skip reserved directories
		}

		directories = append(directories, childDirectory)
	}

	return directories
}
