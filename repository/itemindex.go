// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"io/ioutil"
	"path/filepath"
)

type ItemIndex struct {
	path  string
	items []*Item
}

func NewItemIndex(indexPath string) (*ItemIndex, error) {

	// check if path exists
	if !util.PathExists(indexPath) {
		return nil, fmt.Errorf("The path %q does not exist.", indexPath)
	}

	if isReservedDirectory(indexPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be an index.", indexPath)
	}

	// check if the path is a directory
	if isDirectory, _ := util.IsDirectory(indexPath); !isDirectory {
		indexPath = filepath.Dir(indexPath)
	}

	// create the index
	index := &ItemIndex{
		path:  indexPath,
		items: findAllItems(0, indexPath),
	}

	return index, nil
}

func (itemIndex *ItemIndex) String() string {
	return fmt.Sprintf("%s", itemIndex.path)
}

func (itemIndex *ItemIndex) Path() string {
	return itemIndex.path
}

func (itemIndex *ItemIndex) Directory() string {
	return itemIndex.Path()
}

func (itemIndex *ItemIndex) PathType() string {
	return path.PatherTypeIndex
}

func (itemIndex *ItemIndex) Items() []*Item {
	return itemIndex.items
}

func (itemIndex *ItemIndex) Add(item *Item) {
	itemIndex.items = append(itemIndex.items, item)
}

func findAllItems(level int, itemDirectory string) []*Item {

	items := make([]*Item, 0)

	directoryEntries, err := ioutil.ReadDir(itemDirectory)
	if err != nil {
		fmt.Printf("An error occured while indexing the directory `%v`.\nError: %v\n", itemDirectory, err)
		return nil
	}

	// create a path provider
	pathProviderDirectory := itemDirectory
	if level > 0 {
		pathProviderDirectory = filepath.Dir(itemDirectory)
	}

	pathProvider := path.NewProvider(pathProviderDirectory, false)

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		itemPath := filepath.Join(itemDirectory, element.Name())

		// check if the file a markdown file
		isMarkdown := markdown.IsMarkdownFile(itemPath)
		if !isMarkdown {
			continue
		}

		// search for child items
		childs := getChildItems((level + 1), itemDirectory)

		// create item
		item, err := NewItem(itemPath, level, childs, pathProvider)
		if err != nil {
			fmt.Printf("Skipping item: %s\n", err)
			continue
		}

		// append item to list
		items = append(items, item)

		// item has been found
		directoryContainsItem = true
		break
	}

	// search in sub directories if there is no item in the current folder
	if !directoryContainsItem {

		if virtualItem, err := NewItem(itemDirectory, level, getChildItems((level+1), itemDirectory), pathProvider); err == nil {
			items = append(items, virtualItem)
		} else {
			fmt.Println(err)
		}

	}

	return items
}

func getChildItems(level int, itemDirectory string) []*Item {

	childItems := make([]*Item, 0)

	files, _ := ioutil.ReadDir(itemDirectory)
	for _, folder := range files {

		if !folder.IsDir() {
			continue // skip files
		}

		childItemDirectory := filepath.Join(itemDirectory, folder.Name())
		if isReservedDirectory(childItemDirectory) {
			continue // skip reserved directories
		}

		childsInPath := findAllItems(level, childItemDirectory)
		childItems = append(childItems, childsInPath...)

	}

	return childItems
}
