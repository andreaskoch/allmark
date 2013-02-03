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
	repositoryItems := FindAllRepositoryItems(repositoryPath)
	fmt.Printf("%#v\n", repositoryItems)

	return nil
}

func FindAllRepositoryItems(repositoryPath string) []model.RepositoryItem {

	repositoryItems := make([]model.RepositoryItem, 0, 100)

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v", repositoryPath, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		// check if the file a document
		isNotaRepositoryItem := !strings.EqualFold(strings.ToLower(element.Name()), "notes.md")
		if isNotaRepositoryItem {
			continue
		}

		// search for files
		itemFiles := make([]string, 0, 5)
		filesDirectoryPath := filepath.Join(repositoryPath, "files")
		files, _ := ioutil.ReadDir(filesDirectoryPath)
		for _, element := range files {
			absoluteFilePath := filepath.Join(filesDirectoryPath, element.Name())
			itemFiles = append(itemFiles, absoluteFilePath)
		}

		// search for child items
		childItems := make([]model.RepositoryItem, 0, 100)
		childElements, _ := ioutil.ReadDir(filesDirectoryPath)
		for _, element := range childElements {
			if element.IsDir() {
				folder := filepath.Join(repositoryPath, element.Name())
				childs := FindAllRepositoryItems(folder)
				childItems = append(childItems, childs...)
			}
		}

		// create item and append to list
		item := model.NewRepositoryItem(repositoryPath, itemFiles, childItems)
		repositoryItems = append(repositoryItems, item)

		// item has been found
		directoryContainsItem = true
		break
	}

	// recursive search
	if !directoryContainsItem {
		for _, element := range directoryEntries {

			if element.IsDir() {
				folder := filepath.Join(repositoryPath, element.Name())
				childs := FindAllRepositoryItems(folder)
				repositoryItems = append(repositoryItems, childs...)
			}

		}
	}

	return repositoryItems
}
