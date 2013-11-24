// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"path/filepath"
)

type Repository struct {
	directory string
}

func NewRepository(directory string) (*Repository, error) {

	// check if path exists
	if !fsutil.PathExists(directory) {
		return nil, fmt.Errorf("The path %q does not exist.", directory)
	}

	// check if the supplied path is a file
	if isDirectory, _ := fsutil.IsDirectory(directory); !isDirectory {
		directory = filepath.Dir(directory)
	}

	if isReservedDirectory(directory) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be a root.", directory)
	}

	return &Repository{
		directory: directory,
	}, nil
}

func (itemAccessor *Repository) GetItems() (itemEvents chan *dataaccess.RepositoryEvent, done chan bool) {

	itemEvents = make(chan *dataaccess.RepositoryEvent, 1)
	done = make(chan bool)

	go func() {

		directory := itemAccessor.directory

		// repository directory item
		indexItems(directory, itemEvents)

		done <- true
	}()

	return itemEvents, done
}

// Create a new Item for the specified path.
func indexItems(itemPath string, itemEvents chan *dataaccess.RepositoryEvent) {

	// abort if path does not exist
	if !fsutil.PathExists(itemPath) {
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q does not exist.", itemPath))
		return
	}

	// abort if path is reserved
	if isReservedDirectory(itemPath) {
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemPath))
		return
	}

	// check if its a virtual item or a markdown item
	itemDirectory := filepath.Dir(itemPath)
	if isDirectory, _ := fsutil.IsDirectory(itemPath); isDirectory {

		if found, filepath := findMarkdownFileInDirectory(itemPath); found {

			itemDirectory = itemPath
			itemPath = filepath

		} else {

			// virtual item
			itemDirectory = itemPath

		}

	} else if !isMarkdownFile(itemPath) {

		// the supplied item path does not point to a markdown file
		itemEvents <- dataaccess.NewEvent(nil, fmt.Errorf("%q is not a markdown file.", itemPath))
		return
	}

	// create the file index
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	files := getFiles(filesDirectory)

	// create the item
	item, err := dataaccess.NewItem(itemPath, files)
	itemEvents <- dataaccess.NewEvent(item, err)

	// child items
	childItemDirectories := getChildDirectories(itemDirectory)
	for _, childItemDirectory := range childItemDirectories {
		indexItems(childItemDirectory, itemEvents)
	}
}
