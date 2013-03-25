package indexer

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/parser"
	"github.com/andreaskoch/docs/repository"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func EmptyIndex() *repository.Index {
	return &repository.Index{}
}

func NewIndex(path string) (*repository.Index, error) {

	// check if path is valid
	fileInfo, err := os.Stat(path)
	if err != nil {
		return EmptyIndex(), err
	}

	// check if the path is a directory	
	if !fileInfo.IsDir() {
		return EmptyIndex(), errors.New(fmt.Sprintf("%q is not a directory. Cannot create an index out of a file.", path))
	}

	index := repository.NewIndex(path, findAllItems(path))

	return index, nil
}

func findAllItems(repositoryPath string) []*repository.Item {

	items := make([]*repository.Item, 0, 100)

	directoryEntries, err := ioutil.ReadDir(repositoryPath)
	if err != nil {
		fmt.Printf("An error occured while indexing the repository path `%v`. Error: %v\n", repositoryPath, err)
		return nil
	}

	// item search
	directoryContainsItem := false
	for _, element := range directoryEntries {

		itemPath := filepath.Join(repositoryPath, element.Name())

		// check if the file a markdown file
		isMarkdown := isMarkdownFile(itemPath)
		if !isMarkdown {
			continue
		}

		// search for child items
		childs := getChildItems(repositoryPath)

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
		items = append(items, getChildItems(repositoryPath)...)
	}

	return items
}

func isMarkdownFile(absoluteFilePath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(absoluteFilePath))
	return fileExtension == ".md"
}

func getChildItems(itemPath string) []*repository.Item {

	childItems := make([]*repository.Item, 0, 5)

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
