package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/parser"
	"andyk/docs/util"
	"fmt"
	"os"
	"path/filepath"
)

func GetRenderer(item indexer.Item) (func(), error) {

	// get the lines
	file, err := os.Open(item.Path)
	if err != nil {
		return nil, err
	}
	lines := util.GetLines(file)
	defer file.Close()

	parser := parser.GetParser(lines)
	doc := parser()

	return func() {
		fmt.Println(doc.Title)
	}, nil
}

// Get the filepath of the rendered repository item
func GetRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := filepath.Base(item.Path)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
