// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func EmptyIndex() *repository.ItemIndex {
	return &repository.ItemIndex{}
}

func NewItemIndex(indexDirectory string) (*repository.ItemIndex, error) {

	// check if path is valid
	folderInfo, err := os.Stat(indexDirectory)
	if err != nil {
		return EmptyIndex(), err
	}

	// check if the path is a directory	
	if !folderInfo.IsDir() {
		return EmptyIndex(), errors.New(fmt.Sprintf("%q is not a directory. Cannot create an index out of a file.", indexDirectory))
	}

	return repository.NewItemIndex(indexDirectory, findAllItems(indexDirectory))
}

func findAllItems(itemDirectory string) []*repository.Item {

	items := make([]*repository.Item, 0, 100)

	directoryEntries, err := ioutil.ReadDir(itemDirectory)
	if err != nil {
		fmt.Printf("An error occured while indexing the directory `%v`. Error: %v\n", itemDirectory, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		itemPath := filepath.Join(itemDirectory, element.Name())

		// check if the file a markdown file
		isMarkdown := isMarkdownFile(itemPath)
		if !isMarkdown {
			continue
		}

		// search for child items
		childs := getChildItems(itemDirectory)

		// create item
		item, err := repository.NewItem(itemPath, childs)
		if err != nil {
			fmt.Printf("Skipping item: %s\n", err)
			continue
		}

		// parse item
		if _, err := parser.Parse(item); err != nil {
			fmt.Printf("Could not parse item %q. Error: %s\n", item, err)
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
		items = append(items, getChildItems(itemDirectory)...)
	}

	return items
}

func isMarkdownFile(absoluteFilePath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(absoluteFilePath))
	return fileExtension == ".md"
}

func getChildItems(itemDirectory string) []*repository.Item {

	childItems := make([]*repository.Item, 0, 5)

	files, _ := ioutil.ReadDir(itemDirectory)
	for _, folder := range files {

		if folder.IsDir() {
			childItemDirectory := filepath.Join(itemDirectory, folder.Name())
			childsInPath := findAllItems(childItemDirectory)
			childItems = append(childItems, childsInPath...)
		}

	}

	return childItems
}
