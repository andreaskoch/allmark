// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/markdown"
	util "github.com/andreaskoch/allmark2/common/util/filesystem"
	"github.com/andreaskoch/allmark2/dataaccess"
	"io/ioutil"
	"path/filepath"
)

type ItemAccessor struct {
	directory string
}

func New(directory string) (*ItemAccessor, error) {

	// check if path exists
	if !util.PathExists(directory) {
		return nil, fmt.Errorf("The path %q does not exist.", directory)
	}

	// check if the supplied path is a file
	if isDirectory, _ := util.IsDirectory(directory); !isDirectory {
		directory = filepath.Dir(directory)
	}

	if isReservedDirectory(directory) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be a root.", directory)
	}

	return &ItemAccessor{
		directory: directory,
	}, nil
}

func (itemAccessor *ItemAccessor) GetRootItem() (*dataaccess.Item, error) {

	rootItem, err := newRootItem(itemAccessor.directory)
	if err != nil {
		return nil, err
	}

	return rootItem, nil
}

// Create a new root Item from the specified path.
func newRootItem(path string) (*dataaccess.Item, error) {
	return newItem(path, nil)
}

// Create a new Item for the specified path.
func newItem(itemPath string, parent *dataaccess.Item) (*dataaccess.Item, error) {

	// abort if path does not exist
	if !util.PathExists(itemPath) {
		return nil, fmt.Errorf("The path %q does not exist.", itemPath)
	}

	// abort if path is reserved
	if isReservedDirectory(itemPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an item.", itemPath)
	}

	// check if its a virtual item or a markdown item
	itemDirectory := filepath.Dir(itemPath)
	if isDirectory, _ := util.IsDirectory(itemPath); isDirectory {

		if found, filepath := findMarkdownFileInDirectory(itemPath); found {

			itemDirectory = itemPath
			itemPath = filepath

		} else {

			// virtual item
			itemDirectory = itemPath

		}

	} else if !markdown.IsMarkdownFile(itemPath) {

		// the supplied item path does not point to a markdown file
		return nil, fmt.Errorf("%q is not a markdown file.", itemPath)
	}

	// create the file index
	files := make([]*dataaccess.File, 0)
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)
	if filesRootFolder, err := newRootFolder(filesDirectory); err == nil {
		files = filesRootFolder.Childs()
	}

	// create the item
	item, err := dataaccess.NewItem(itemPath, parent, files)
	if err != nil {
		return nil, err
	}

	// append childs
	item.SetChilds(getChildItems(itemDirectory))

	return item, nil
}

func getChildItems(directory string) []*dataaccess.Item {

	newChildList := make([]*dataaccess.Item, 0)

	// add new childs
	childItemDirectories := getChildDirectories(directory)
	for _, childItemDirectory := range childItemDirectories {
		if child, err := newRootItem(childItemDirectory); err == nil {
			newChildList = append(newChildList, child) // append new child
		} else {
			fmt.Printf("Could not create a item for folder %q. Error: %s\n", childItemDirectory, err)
		}
	}

	return newChildList
}

func newRootFolder(path string) (*dataaccess.File, error) {
	return newFile(nil, path)
}

func newFile(parent *dataaccess.File, path string) (*dataaccess.File, error) {

	// check if the path exists
	if exists := util.PathExists(path); !exists {
		return nil, fmt.Errorf("The path %q does not exists.", path)
	}

	// check if the path is a directory or a file
	isDir, err := util.IsDirectory(path)
	if err != nil {
		return nil, err
	}

	// create the file
	file, err := dataaccess.NewFile(path, parent)
	if err != nil {
		return nil, err
	}

	// append childs
	if isDir {
		file.SetChilds(getChildFiles(path))
	}

	return file, nil
}

func getChildFiles(directory string) []*dataaccess.File {

	childs := make([]*dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		// append new file
		path := filepath.Join(directory, directoryEntry.Name())
		file, err := newRootFolder(path)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", path, err)
		}

		childs = append(childs, file)
	}

	return childs
}
