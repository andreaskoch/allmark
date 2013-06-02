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
	"time"
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
	Parent *Item
	Childs []*Item

	newChild     chan *Item
	deletedChild chan *Item

	directory string
	path      string
	isVirtual bool

	pathProvider *path.Provider

	changeFuncs []func(item *Item)
}

func newItem(parent *Item, itemPath string, level int, newItem chan *Item, deletedItem chan *Item) (item *Item, err error) {

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

		Parent: parent,
		Level: level,
		Files: fileIndex,

		newChild:     newItem,
		deletedChild: deletedItem,

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

					fmt.Printf("Item %q has been modified\n", item)
					item.Modified <- true

					// update parent
					if item.Parent != nil {
						item.Parent.Modified <- true
					}

				case <-fileWatcher.Moved:

					fmt.Printf("Item %q has been moved\n", item)
					item.Moved <- true
				}

			}
		}()
	}

	// look for changes in the item directory
	go func() {
		var skipFunc = func(path string) bool {
			isItem := path == item.path
			isReserved := isReservedDirectory(path)
			return isItem || isReserved
		}

		folderWatcher := watcher.NewFolderWatcher(itemDirectory, false, skipFunc).Start()

		for folderWatcher.IsRunning() {

			select {
			case <-folderWatcher.Change:
				fmt.Printf("Updating the childs of item %q\n", item)
				item.updateChilds()

				go func() {
					time.Sleep(time.Second * 3)
					item.Modified <- true
				}()
			}

		}
	}()

	// look for file changes
	go func() {
		for {
			select {
			case <-item.Files.Changed:
				fmt.Printf("Files of item %q changed\n", item)
				item.Modified <- true
			case <-item.Files.Stopped:
				fmt.Printf("File index watcher for item %q was stopped\n", item)
				break
			}
		}

		fmt.Println("File watcher stopped")
	}()

	return item, nil
}

func (item *Item) OnChange(name string, expr func(i *Item)) {
	item.changeFuncs = append(item.changeFuncs, expr)
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.PathProvider().GetWebRoute(item))
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

	// check if the child list needs initialization
	if item.Childs == nil {
		item.Childs = make([]*Item, 0, len(childItemDirectories))
	}

	// add new childs
	for _, childItemDirectory := range childItemDirectories {

		if item.HasChild(childItemDirectory) {
			continue // Child already present
		}

		if child, err := newItem(item, childItemDirectory, item.Level+1, item.newChild, item.deletedChild); err == nil {

			// inform others about the new child
			go func() {
				item.newChild <- child
			}()

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

		} else {

			// Todo: stop all go routines on this child

			// inform others about the removed child
			go func() {
				item.deletedChild <- child
			}()

		}
	}

	// assign new child list
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

		// append directory
		directories = append(directories, childDirectory)
	}

	return directories
}
