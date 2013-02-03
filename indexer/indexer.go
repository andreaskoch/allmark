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

func Index(repositoryPath string) map[int]model.Document {

	// check if the supplied repository path is set
	if strings.Trim(repositoryPath, " ") == "" {
		fmt.Print("The repository path cannot be null or empty.")
		return nil
	}

	// check if the supplied repository path exists
	if _, err := os.Stat(repositoryPath); err != nil {
		switch {
		case os.IsNotExist(err):
			fmt.Printf("The supplied repository path `%v` does not exist.", repositoryPath)
		default:
			fmt.Printf("An error occured while trying to access the supplied repository path `%v`.", repositoryPath)
		}

		return nil
	}

	// get all repository items in the supplied repository path
	repositoryItems := make([]*model.RepositoryItem, 100)
	FindAllRepositoryItems(repositoryPath, repositoryItems)
	fmt.Printf("%v", repositoryItems)

	return nil
}

func FindAllRepositoryItems(repositoryPath string, repositoryItems []*model.RepositoryItem) {

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v", repositoryPath, err)
		return
	}

	for _, element := range directoryEntries {

		if element.IsDir() {
			//fmt.Printf("Element `%v` is a directory. Recurse.\n", element.Name())
			FindAllRepositoryItems(filepath.Join(repositoryPath, element.Name()), repositoryItems)
		}

		// check if the file a document
		isRepositoryItem := strings.EqualFold(strings.ToLower(element.Name()), "notes.md")
		if !isRepositoryItem {
			continue
		}

		newItem := model.NewRepositoryItem(repositoryPath)
		numberOfItems := getNumberOfItemsInArray(repositoryItems)
		repositoryItems[numberOfItems] = newItem
	}
}

func getNumberOfItemsInArray(arr []*model.RepositoryItem) int {
	if arr == nil || len(arr) == 0 {
		return 0
	}

	for index, element := range arr {
		if element == nil {
			return index
		}
	}

	return len(arr)
}
