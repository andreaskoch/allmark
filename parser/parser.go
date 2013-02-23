package parser

import (
	"andyk/docs/indexer"
	"andyk/docs/util"
	"os"
)

func ParseItem(item indexer.Item) (Document, error) {

	// open the file
	file, err := os.Open(item.Path)
	if err != nil {
		return Document{}, err
	}

	defer file.Close()

	// get the lines
	lines := util.GetLines(file)

	// parse the document
	document := CreateDocument(lines)

	return document, nil
}
