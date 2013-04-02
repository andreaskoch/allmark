package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/mapper"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/templates"
	"github.com/andreaskoch/docs/view"
	"os"
	"text/template"
)

func RenderRepositories(repositoryPaths []string) []*repository.Index {
	numberOfRepositories := len(repositoryPaths)
	indizes := make([]*repository.Index, numberOfRepositories, numberOfRepositories)

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

func renderIndex(index *repository.Index) *repository.Index {

	renderFunc := func(item *repository.Item) *repository.Item {
		indexDirectory := index.Path

		pathProviderFunc := func(item *repository.Item) string {
			return item.GetRelativePath(indexDirectory)
		}

		return renderItem(item, pathProviderFunc)
	}

	index.Walk(func(item *repository.Item) {

		// render the item
		item.Render(renderFunc)

		// render the item again if it changes
		item.RegisterOnChangeCallback("RenderOnChange", func(i *repository.Item) {

			fmt.Printf("Item %q changed", item)
			i.Render(renderFunc)
		})
	})

	return index
}

func renderItem(item *repository.Item, pathProviderFunc func(item *repository.Item) string) *repository.Item {

	fmt.Printf("Rendering item %q\n", item.Path)

	fmt.Println("Reindexing files")
	item.IndexFiles()

	// get a template
	templateText, err := templates.GetTemplate(item)
	if err != nil {
		fmt.Println(err)
		return item
	}

	// get a viewmodel mapper
	mapperFunc, err := mapper.GetMapper(item, pathProviderFunc)
	if err != nil {
		fmt.Println(err)
		return item
	}

	// create the viewmodel
	viewModel := mapperFunc(item)

	// render the template
	render(item, templateText, viewModel)

	return item
}

func render(item *repository.Item, templateText string, viewModel view.Model) (*repository.Item, error) {
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
