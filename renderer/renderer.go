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

func renderItem(item indexer.Item) interface{} {

	_, err := parser.Parse(&item)
	if err != nil {
		fmt.Printf("Could not parse item \"%v\": %v\n", item.Path, err)
		return nil
	}

	switch item.Type {
	case indexer.DocumentItemType:
		{
			document := mappers.GetDocument(item)
			render(item, templates.DocumentTemplate, document)
			return document
		}

	case indexer.RepositoryItemType:
		{
			repository := mappers.GetRepository(item)
			render(item, templates.RepositoryTemplate, repository)
			return repository
		}
	}

	return nil
}

func render(item indexer.Item, templateText string, viewModel interface{}) {
	file, err := os.Create(item.RenderedPath)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	template := template.New(item.Type)
	template.Parse(templateText)
	template.Execute(writer, viewModel)
}
