package parser

import (
	"andyk/docs/indexer"
	"andyk/docs/util"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ParsedItemElement struct {
	Name  string
	Value string
}

func NewParsedItemElement(name string, value string) ParsedItemElement {
	return ParsedItemElement{
		Name:  name,
		Value: value,
	}
}

type ParsedItem struct {
	Elements []ParsedItemElement
	Type     string
	MetaData MetaData
}

func (parsedItem *ParsedItem) AddElement(name string, value string) {

	element := NewParsedItemElement(name, value)

	if parsedItem.Elements == nil {
		parsedItem.Elements = make([]ParsedItemElement, 1, 1)
		parsedItem.Elements[0] = element
		return
	}

	parsedItem.Elements = append(parsedItem.Elements, element)

}

type Parser interface {
	Parse(lines []string) (ParsedItem, error)
}

func ParseItem(item indexer.Item) (ParsedItem, error) {

	// open the file
	file, err := os.Open(item.Path)
	if err != nil {
		return ParsedItem{}, err
	}

	defer file.Close()

	// get the lines
	lines := util.GetLines(file)

	// define the patterns
	documentStructure := NewDocumentStructure()

	// get the meta data
	metaDataParser := NewMetaDataParser(documentStructure)
	metaData, metaDataLocation, lines := metaDataParser.Parse(lines)
	if !metaDataLocation.Found {

		// infer type from file name
		metaData.ItemType = getItemTypeFromFilename(item.GetFilename())

	}

	// parse by type
	itemType := strings.TrimSpace(strings.ToLower(metaData.ItemType))
	switch itemType {
	case "document":
		{
			parser := NewDocumentParser(documentStructure)
			return parser.Parse(lines, metaData)
		}
	}

	return ParsedItem{}, errors.New(fmt.Sprintf("Items of type \"%v\" cannot be parsed.", itemType))
}

func getItemTypeFromFilename(filename string) string {

	lowercaseFilename := strings.ToLower(filename)

	switch lowercaseFilename {
	case "repository.md":
		return "repository"

	case "document.md":
		return "document"

	case "location.md":
		return "location"

	case "comment.md":
		return "comment"

	case "message.md":
		return "message"
	}

	return "unknown"
}
