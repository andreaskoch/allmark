// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"andyk/docs/repository"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Index(repositoryPath string) repository.Index {

	// check if the supplied repository path is set
	if strings.Trim(repositoryPath, " ") == "" {
		panic("The repository path cannot be null or empty.")
	}

	// check if the supplied repository path exists
	if _, err := os.Stat(repositoryPath); err != nil {

		switch {
		case os.IsNotExist(err):
			panic(fmt.Sprintf("The supplied repository path `%v` does not exist.", repositoryPath))
		default:
			panic(fmt.Sprintf("An error occured while trying to access the supplied repository path `%v`.", repositoryPath))
		}
	}

	// get all repository items in the supplied repository path
	items := findAllItems(repositoryPath)

	index := repository.NewIndex(repositoryPath, items)

	return index
}

func findAllItems(repositoryPath string) []repository.Item {

	items := make([]repository.Item, 0, 100)

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v\n", repositoryPath, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		// check if the file a document
		isItem, itemType := fileIsItem(element.Name())
		if !isItem {
			continue
		}

		// search for files
		files := getFiles(repositoryPath)

		// search for child items
		childs := getChildItems(repositoryPath)

		// create item and append to list
		itemPath := filepath.Join(repositoryPath, element.Name())
		item := repository.NewItem(itemType, itemPath, files, childs)
		items = append(items, item)

		// item has been found
		directoryContainsItem = true
		break
	}

	// search in sub directories if there is no item in the current folder
	if !directoryContainsItem {
		items = append(items, getChildItems(repositoryPath)...)
	}

	return items
}

func fileIsItem(filename string) (bool, string) {

	lowercaseFilename := strings.ToLower(filename)

	switch lowercaseFilename {
	case "repository.md":
		return true, "repository"

	case "document.md":
		return true, "document"

	case "location.md":
		return true, "location"

	case "comment.md":
		return true, "comment"

	case "message.md":
		return true, "message"
	}

	return false, "unknown"
}

func getChildItems(itemPath string) []repository.Item {

	childItems := make([]repository.Item, 0, 5)

	files, _ := ioutil.ReadDir(itemPath)
	for _, element := range files {

		if element.IsDir() {
			path := filepath.Join(itemPath, element.Name())
			childsInPath := findAllItems(path)
			childItems = append(childItems, childsInPath...)
		}

	}

	return childItems
}

func getFiles(itemPath string) []repository.File {

	itemFiles := make([]repository.File, 0, 5)
	filesDirectoryEntries, _ := ioutil.ReadDir(filepath.Join(itemPath, "files"))

	for _, file := range filesDirectoryEntries {
		absoluteFilePath := filepath.Join(itemPath, file.Name())
		repositoryFile := repository.NewFile(absoluteFilePath)

		itemFiles = append(itemFiles, repositoryFile)
	}

	return itemFiles
}
