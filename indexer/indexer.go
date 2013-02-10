// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"andyk/docs/model"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Index(repositoryPath string) model.RepositoryIndex {

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
	repositoryItems := findAllRepositoryItems(repositoryPath)

	index := model.NewRepositoryIndex(repositoryPath, repositoryItems)

	return index
}

func findAllRepositoryItems(repositoryPath string) []model.RepositoryItem {

	repositoryItems := make([]model.RepositoryItem, 0, 100)

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v\n", repositoryPath, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		// check if the file a document
		isRepositoryItem, itemType := fileIsRepositoryItem(element.Name())
		if !isRepositoryItem {
			continue
		}

		// search for files
		files := getFiles(repositoryPath)

		// search for child items
		childs := getChildItems(repositoryPath)

		// create item and append to list
		itemPath := filepath.Join(repositoryPath, element.Name())
		item := model.NewRepositoryItem(itemType, itemPath, files, childs)
		repositoryItems = append(repositoryItems, item)

		// item has been found
		directoryContainsItem = true
		break
	}

	// search in sub directories if there is no item in the current folder
	if !directoryContainsItem {
		repositoryItems = append(repositoryItems, getChildItems(repositoryPath)...)
	}

	return repositoryItems
}

func fileIsRepositoryItem(filename string) (bool, string) {

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

func getChildItems(repositoryItemPath string) []model.RepositoryItem {

	childItems := make([]model.RepositoryItem, 0, 5)

	files, _ := ioutil.ReadDir(repositoryItemPath)
	for _, element := range files {

		if element.IsDir() {
			path := filepath.Join(repositoryItemPath, element.Name())
			childsInPath := findAllRepositoryItems(path)
			childItems = append(childItems, childsInPath...)
		}

	}

	return childItems
}

func getFiles(repositoryItemPath string) []model.RepositoryItemFile {

	itemFiles := make([]model.RepositoryItemFile, 0, 5)
	filesDirectoryEntries, _ := ioutil.ReadDir(filepath.Join(repositoryItemPath, "files"))

	for _, file := range filesDirectoryEntries {
		absoluteFilePath := filepath.Join(repositoryItemPath, file.Name())
		repositoryItemFile := model.NewRepositoryItemFile(absoluteFilePath)

		itemFiles = append(itemFiles, repositoryItemFile)
	}

	return itemFiles
}
