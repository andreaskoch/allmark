package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/mapper"
	"github.com/andreaskoch/docs/parser"
	"github.com/andreaskoch/docs/templates"
	"os"
	"text/template"
)

func RenderRepositories(repositoryPaths []string) []*indexer.Index {
	numberOfRepositories := len(repositoryPaths)
	indizes := make([]*indexer.Index, numberOfRepositories, numberOfRepositories)

	for indexNumber, repositoryPath := range repositoryPaths {
		index, err := indexer.NewIndex(repositoryPath)
		if err != nil {
			fmt.Printf("Cannot create an index for folder %q. Error: %v", repositoryPath, err)
			continue
		}

		indizes[indexNumber] = renderIndex(index)
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
		return item
	}

	fmt.Println("Reindexing files")
	item.IndexFiles()

	// get a template
	templateText, err := templates.GetTemplate(item)
	if err != nil {
		fmt.Println(err)
		return item
	}

	// get a viewmodel mapper
	mapperFunc, err := mapper.GetMapper(item)
	if err != nil {
		fmt.Println(err)
		return item
	}

	// create the viewmodel
	viewModel := mapperFunc(item, nil)

	// render the template
	render(item, templateText, viewModel)

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
