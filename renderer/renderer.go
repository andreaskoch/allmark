// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/indexer"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/renderer/html"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"github.com/andreaskoch/allmark/watcher"
	"os"
	"text/template"
)

func RenderRepository(repositoryPath string) *repository.ItemIndex {
	itemIndex, err := indexer.NewItemIndex(repositoryPath)
	if err != nil {
		fmt.Printf("Cannot create an item index for folder %q. Error: %v", repositoryPath, err)
		return nil
	}

	return renderIndex(itemIndex)
}

func renderIndex(itemIndex *repository.ItemIndex) *repository.ItemIndex {

	repositoryPath := itemIndex.Path()
	for _, item := range itemIndex.Items() {
		renderItem(repositoryPath, item)
	}

	return itemIndex
}

func renderItem(repositoryPath string, item *repository.Item) (*repository.Item, error) {

	fmt.Printf("RENDERING: %s\n", item)

	// render child items first
	for _, child := range item.Childs() {

		// attach change listener
		child.OnChange("Throw Item Events on Child Item change", func(event *watcher.WatchEvent) {
			item.Throw(event)
		})

		renderItem(repositoryPath, child)
	}

	// parse the item
	if _, parseError := parser.Parse(item); parseError != nil {
		return item, fmt.Errorf("Cannot render the item %q, because it could not be parsed.\nError: %s\n", item, parseError)
	}

	// attach change listener
	item.OnChange("Render item on change", func(event *watcher.WatchEvent) {
		renderItem(repositoryPath, item)
	})

	// get a template
	templateText, err := templates.GetTemplate(item)
	if err != nil {
		return item, err
	}

	// create a path provider
	pathProvider := path.NewProvider(repositoryPath)

	// get a viewmodel mapper
	mapperFunc, err := mapper.GetMapper(pathProvider, html.NewConverter, item.Type)
	if err != nil {
		return item, err
	}

	// create the viewmodel
	viewModel := mapperFunc(item)

	// render the template
	render(item, templateText, viewModel)

	return item, nil
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
