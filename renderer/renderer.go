package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/mappers"
	"andyk/docs/parser"
	"andyk/docs/templates"
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Render(repositoryPaths []string) {
	for _, repositoryPath := range repositoryPaths {
		index := indexer.GetIndex(repositoryPath)
		renderIndex(index)
	}
}

func renderIndex(index indexer.Index) {

	for _, item := range index.Items {
		renderItem(item)
	}

}

func renderItem(item indexer.Item) {

	parsedItem, err := parser.Parse(item)
	if err != nil {
		log.Printf("Could not parse item \"%v\". Error: %v", item.Path, err)
		return
	}

	renderedItemFilePath := getRenderedItemPath(item)

	switch parsedItem.MetaData.ItemType {
	case parser.DocumentItemType:
		{
			file, err := os.Create(renderedItemFilePath)
			if err != nil {
				panic(err)
			}
			writer := bufio.NewWriter(file)

			defer func() {
				writer.Flush()
				file.Close()
			}()

			document := mappers.GetDocument(parsedItem)
			template := template.New(parser.DocumentItemType)
			template.Parse(templates.DocumentTemplate)
			template.Execute(writer, document)
		}
	}
}

// Get the filepath of the rendered repository item
func getRenderedItemPath(item indexer.Item) string {
	itemDirectory := filepath.Dir(item.Path)
	itemName := strings.Replace(filepath.Base(item.Path), filepath.Ext(item.Path), "", 1)

	renderedFilePath := filepath.Join(itemDirectory, itemName+".html")
	return renderedFilePath
}
