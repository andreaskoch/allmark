package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/mapper"
	"github.com/andreaskoch/docs/parser"
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/templates"
	"github.com/andreaskoch/docs/view"
	"os"
	"text/template"
)

func RenderRepositories(repositoryPaths []string) []*repository.ItemIndex {
	numberOfRepositories := len(repositoryPaths)
	indizes := make([]*repository.ItemIndex, numberOfRepositories, numberOfRepositories)

	for indexNumber, repositoryPath := range repositoryPaths {
		itemIndex, err := indexer.NewItemIndex(repositoryPath)
		if err != nil {
			fmt.Printf("Cannot create an item index for folder %q. Error: %v", repositoryPath, err)
			continue
		}

		indizes[indexNumber] = renderIndex(itemIndex)
	}

	return indizes
}

func renderIndex(itemIndex *repository.ItemIndex) *repository.ItemIndex {

	itemIndex.Walk(func(item *repository.Item) {

		// render the item
		renderItem(itemIndex.Path(), item)

		// render the item again if it changes
		item.RegisterOnChangeCallback("RenderOnChange", func(i *repository.Item) {

			fmt.Printf("Item %q changed", item)
			if _, parseError := parser.Parse(item); parseError == nil {
				renderItem(itemIndex.Path(), item)
			} else {
				fmt.Printf("Cannot render the item %q, because it could not be parsed. Error: %s", item, parseError)
			}

		})
	})

	return itemIndex
}

func renderItem(repositoryPath string, item *repository.Item) *repository.Item {

	fmt.Printf("Rendering item %q\n", item)

	// get a template
	templateText, err := templates.GetTemplate(item)
	if err != nil {
		fmt.Println(err)
		return item
	}

	// create a path provider
	pathProvider := path.NewProvider(repositoryPath)

	// get a viewmodel mapper
	mapperFunc, err := mapper.GetMapper(pathProvider, item)
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

	targetPath := path.GetRenderTargetPath(item)
	file, err := os.Create(targetPath)
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
