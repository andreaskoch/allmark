package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/mappers"
	"andyk/docs/parser"
	"andyk/docs/templates"
	"bufio"
	"fmt"
	"os"
	"text/template"
)

func Render(repositoryPaths []string) {
	for _, repositoryPath := range repositoryPaths {
		index := indexer.GetIndex(repositoryPath)
		RenderIndex(index)
	}
}

func RenderIndex(index indexer.Index) {

	index.Walk(func(item indexer.Item) {
		renderItem(item)
	})

}

func renderItem(item indexer.Item) {

	parsedItem, err := parser.Parse(item)
	if err != nil {
		fmt.Printf("Could not parse item \"%v\": %v\n", item.Path, err)
		return
	}

	switch parsedItem.Item.Type {
	case indexer.DocumentItemType:
		{
			file, err := os.Create(item.RenderedPath)
			if err != nil {
				panic(err)
			}
			writer := bufio.NewWriter(file)

			defer func() {
				writer.Flush()
				file.Close()
			}()

			document := mappers.GetDocument(parsedItem)
			template := template.New(indexer.DocumentItemType)
			template.Parse(templates.DocumentTemplate)
			template.Execute(writer, document)
		}
	}
}
