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

func Render(repositoryPaths []string) []indexer.Index {
	indizes := make([]indexer.Index, len(repositoryPaths), len(repositoryPaths))

	for _, repositoryPath := range repositoryPaths {
		index := indexer.GetIndex(repositoryPath)
		indizes = append(indizes, renderIndex(index))
	}

	return indizes
}

func renderIndex(index indexer.Index) indexer.Index {

	index.Walk(func(item indexer.Item) {
		renderItem(item)
	})

	return index
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

	case indexer.MessageItemType:
		{
			message := mappers.GetMessage(item)
			render(item, templates.MessageTemplate, message)
			return message
		}

	case indexer.CollectionItemType:
		{
			makeSureChildItemsAreParsed := func(item *indexer.Item) {
				parser.Parse(item)
			}

			collection := mappers.GetCollection(item, makeSureChildItemsAreParsed)
			render(item, templates.CollectionTemplate, collection)
			return collection
		}

	case indexer.RepositoryItemType:
		{
			makeSureChildItemsAreParsed := func(item *indexer.Item) {
				parser.Parse(item)
			}

			repository := mappers.GetRepository(item, makeSureChildItemsAreParsed)
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
