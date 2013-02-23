package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/parser"
	"fmt"
	"log"
	"path/filepath"
)

func RenderItem(item indexer.Item) {

	parsedItem, err := parser.ParseItem(item)
	if err != nil {
		log.Printf("Could not parse item \"%v\". Error: %v", item.Path, err)
		return
	}

	fmt.Printf("%#v", parsedItem)
}

// Get the filepath of the rendered repository item
func GetRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := filepath.Base(item.Path)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
