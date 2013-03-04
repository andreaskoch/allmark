package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/mappers"
	"github.com/andreaskoch/docs/parser"
	"github.com/andreaskoch/docs/templates"
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

	_, err := parser.Parse(&item)
	if err != nil {
		fmt.Printf("Could not parse item \"%v\": %v\n", item.Path, err)
		return
	}

	switch item.Type {
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

			document := mappers.GetDocument(item)
			template := template.New(indexer.DocumentItemType)
			template.Parse(templates.DocumentTemplate)
			template.Execute(writer, document)
		}
	}
}
