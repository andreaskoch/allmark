package repository

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

func GetRenderer(repositoryItem *Item) func() {

	return func() {

	}

}

func GetParser(repositoryItem *Item) func() {

	return func() {

	}

}

// Get the hash code of the rendered item
func GetRenderedItemHash(item Item) string {
	renderedItemPath := GetRenderedItemPath(item)

	file, err := os.Open(renderedItemPath)
	if err != nil {
		// file does not exist or cannot be accessed
		return ""
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	firstLineBytes, _ := fileReader.ReadBytes('\n')
	if firstLineBytes == nil {
		// first line cannot be read
		return ""
	}

	// extract hash from line
	hashCodeRegexp := regexp.MustCompile("<!-- (\\w+) -->")
	matches := hashCodeRegexp.FindStringSubmatch(string(firstLineBytes))
	if len(matches) != 2 {
		return ""
	}

	extractedHashcode := matches[1]

	return string(extractedHashcode)
}

// Get the filepath of the rendered repository item
func GetRenderedItemPath(item Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := filepath.Base(item.Path)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
