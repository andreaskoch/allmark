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

	Title            string
	Description      string
	RawContent       []string
	ConvertedContent string
	MetaData         MetaData

	AbsolutePath string
	RelativePath string

	ChildsReady  chan bool
	newChild     chan *Item
	deletedChild chan *Item

	directory string
	path      string
	isVirtual bool

	rootPathProvider     *path.Provider
	relativePathProvider *path.Provider
	filePathProvider     *path.Provider
}

func newItem(rootPathProvider *path.Provider, parent *Item, itemPath string, level int, newItem chan *Item, deletedItem chan *Item) (*Item, error) {

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

	// create a relative path provider
	pathProviderDirectory := itemDirectory
	if level > 0 {
		pathProviderDirectory = filepath.Dir(itemDirectory)
	}

	relativePathProvider := rootPathProvider.New(pathProviderDirectory)

	// create a file path provider
	filePathProvider := rootPathProvider.New(itemDirectory)

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, FilesDirectoryName)
	fileIndex, err := newFileIndex(rootPathProvider, filesDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a file index for folder %q.\nError: %s\n", filesDirectory, err)
	}

	// create the item
	item := &Item{

		ChildsReady: make(chan bool),
		Modified:    make(chan bool),
		Moved:       make(chan bool),

		Parent: parent,
		Level:  level,
		Files:  fileIndex,

		newChild:     newItem,
		deletedChild: deletedItem,

		rootPathProvider:     rootPathProvider,
		relativePathProvider: relativePathProvider,
		filePathProvider:     filePathProvider,

		directory: itemDirectory,
		path:      itemPath,
		isVirtual: isVirtualItem,
	}

	// assign paths
	item.RelativePath = item.RelativePathProvider().GetWebRoute(item)
	item.AbsolutePath = item.RootPathProvider().GetWebRoute(item)

	// find childs
	item.updateChilds()

	// look for changes to the markdown file (if the item is not virtual)
	if !item.isVirtual {
		item.startFileWatcher()
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

				if item.isVirtual {
					if found, filepath := findMarkdownFileInDirectory(itemDirectory); found {

						fmt.Printf("Converting the virtual item %q into a physical item\n", item)

						// make the item physical
						item.path = filepath
						item.isVirtual = false
						item.startFileWatcher()
					}
				}

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
	}()

	return item, nil
}

func (item *Item) startFileWatcher() (started bool, itemWatcher *watcher.FileWatcher) {
	if isFile, _ := util.IsFile(item.path); !isFile {
		return false, nil
	}

	fileWatcher := watcher.NewFileWatcher(item.path).Start()

	go func() {

		for fileWatcher.IsRunning() {

			select {
			case <-fileWatcher.Modified:

				fmt.Printf("Item %q has been modified\n", item)
				item.Modified <- true

			case <-fileWatcher.Moved:

				fmt.Printf("Item %q has been moved\n", item)
				item.Moved <- true
			}

		}
	}()

	return true, fileWatcher
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.RootPathProvider().GetWebRoute(item))
}

func (item *Item) Less(otherItem *Item) bool {
	return item.MetaData.Date.Before(otherItem.MetaData.Date)
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

func (item *Item) RootPathProvider() *path.Provider {
	return item.rootPathProvider
}

func (item *Item) RelativePathProvider() *path.Provider {
	return item.relativePathProvider
}

func (item *Item) FilePathProvider() *path.Provider {
	return item.filePathProvider
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

		if child, err := newItem(item.RootPathProvider(), item, childItemDirectory, item.Level+1, item.newChild, item.deletedChild); err == nil {

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
	go func() {
		item.ChildsReady <- true
	}()

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
