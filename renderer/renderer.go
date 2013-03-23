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

func Render(repositoryPaths []string) []*indexer.Index {
	indizes := make([]*indexer.Index, len(repositoryPaths), len(repositoryPaths))

	for _, repositoryPath := range repositoryPaths {
		index, err := indexer.NewIndex(repositoryPath)
		if err != nil {
			fmt.Printf("Cannot create an index for folder %q. Error: %v", repositoryPath, err)
			continue
		}

		indizes = append(indizes, renderIndex(index))
	}

	return indizes
}

func renderIndex(index *indexer.Index) *indexer.Index {

	index.Walk(func(item *indexer.Item) {

		// render the item
		item.Render(renderItem)

		// render the item again if it changes
		item.RegisterOnChangeCallback("RenderOnChange", func(i *indexer.Item) {

			fmt.Printf("Item %q changed", item)
			i.Render(renderItem)
		})
	})

	return index
}

func renderItem(item *indexer.Item) *indexer.Item {

	fmt.Printf("Rendering item %q\n", item.Path)

	_, err := parser.Parse(item)
	if err != nil {
		fmt.Printf("Could not parse item \"%v\": %v\n", item.Path, err)
		return nil
	}

	fmt.Println("Reindexing files")
	item.IndexFiles()

	switch item.Type {
	case indexer.DocumentItemType:
		{
			document := mappers.GetDocument(*item)
			render(item, templates.DocumentTemplate, document)
		}

	case indexer.MessageItemType:
		{
			message := mappers.GetMessage(*item)
			render(item, templates.MessageTemplate, message)
		}

	case indexer.CollectionItemType:
		{
			makeSureChildItemsAreParsed := func(item *indexer.Item) {
				parser.Parse(item)
			}

			collection := mappers.GetCollection(*item, makeSureChildItemsAreParsed)
			render(item, templates.CollectionTemplate, collection)
		}

	case indexer.RepositoryItemType:
		{
			makeSureChildItemsAreParsed := func(item *indexer.Item) {
				parser.Parse(item)
			}

			repository := mappers.GetRepository(*item, makeSureChildItemsAreParsed)
			render(item, templates.RepositoryTemplate, repository)
		}
	}

	return item
}

func render(item *indexer.Item, templateText string, viewModel interface{}) (*indexer.Item, error) {
	file, err := os.Create(item.RenderedPath)
	if err != nil {
		return item, err
	}

	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	template := template.New(item.Type)
	template.Parse(templateText)
	template.Execute(writer, viewModel)

	return item, nil
}
