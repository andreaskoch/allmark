package repository

import (
	"path/filepath"
)

func GetRenderer(repositoryItem *Item) func() {

	return func() {

	}

}

func GetParser(repositoryItem *Item) func() {

	return func() {

	}

}

// Get the filepath of the rendered repository item
func GetRenderedItemPath(item Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := filepath.Base(item.Path)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
