package renderer

import (
	"andyk/docs/indexer"
	"andyk/docs/mappers"
	"andyk/docs/parser"
	"andyk/docs/templates"
	"bufio"
	"log"
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

	switch parsedItem.MetaData.ItemType {
	case parser.DocumentItemType:
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
			template := template.New(parser.DocumentItemType)
			template.Parse(templates.DocumentTemplate)
			template.Execute(writer, document)
		}
	}
}
