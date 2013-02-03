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

	for index, repositoryItem := range repositoryItems {
		fmt.Println(index)
		fmt.Println(repositoryItem.Path)

		fmt.Println("Files:")
		for _, file := range repositoryItem.Files {
			fmt.Println(file)
		}
		fmt.Println()

		fmt.Println("Child elements:")
		for _, child := range repositoryItem.ChildItems {
			fmt.Println(child.Path)
		}

		fmt.Println()
	}

	return nil
}

func FindAllRepositoryItems(repositoryPath string) []model.RepositoryItem {

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
		isNotaRepositoryItem := !strings.EqualFold(strings.ToLower(element.Name()), "notes.md")
		if isNotaRepositoryItem {
			continue
		}

		// search for files
		files := GetFiles(repositoryPath)

		// search for child items
		childs := GetChildItems(repositoryPath)

		// create item and append to list
		item := model.NewRepositoryItem(repositoryPath, files, childs)
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

func GetChildItems(repositoryItemPath string) []model.RepositoryItem {

	childItems := make([]model.RepositoryItem, 0, 5)

	files, _ := ioutil.ReadDir(repositoryItemPath)
	for _, element := range files {

		if element.IsDir() {
			path := filepath.Join(repositoryItemPath, element.Name())
			childsInPath := FindAllRepositoryItems(path)
			childItems = append(childItems, childsInPath...)
		}

	}

	return childItems
}

func GetFiles(repositoryItemPath string) []string {

	itemFiles := make([]string, 0, 1)
	filesDirectoryEntries, _ := ioutil.ReadDir(filepath.Join(repositoryItemPath, "files"))

	for _, file := range filesDirectoryEntries {
		absoluteFilePath := filepath.Join(repositoryItemPath, file.Name())
		itemFiles = append(itemFiles, absoluteFilePath)
	}

	return itemFiles
}
