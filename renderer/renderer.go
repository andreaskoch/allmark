package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/parser"
	"fmt"
	"path/filepath"
)

func GetRenderer(item indexer.Item) func() {

	parser := parser.GetParser(make([]string, 0))
	doc := parser()

	return func() {
		fmt.Println(doc.Title)
	}

}

// Get the filepath of the rendered repository item
func GetRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := filepath.Base(item.Path)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
