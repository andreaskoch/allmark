// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"github.com/andreaskoch/allmark/watcher"
	"os"
	"text/template"
)

func RenderRepository(repositoryPath string) *repository.ItemIndex {
	itemIndex, err := repository.NewItemIndex(repositoryPath)
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

	// render child items first
	for _, child := range item.Childs() {

		// attach change listener
		child.OnChange("Throw Item Events on Child Item change", func(event *watcher.WatchEvent) {
			item.Throw(event)
		})

		renderItem(repositoryPath, child)
	}

	// attach change listener
	item.OnChange("Render item on change", func(event *watcher.WatchEvent) {
		renderItem(repositoryPath, item)
	})

	// create a path provider
	pathProvider := path.NewProvider(repositoryPath)

	// create the viewmodel
	viewModel := mapper.Map(item, pathProvider, "html")

	// get a template
	templateText := templates.GetTemplate(viewModel.Type)

	// render the template
	render(item, pathProvider, templateText, viewModel)

	return item, nil
}

func render(item *repository.Item, pathProvider *path.Provider, templateText string, viewModel view.Model) (*repository.Item, error) {

	targetPath := pathProvider.GetRenderTargetPath(item)
	file, err := os.Create(targetPath)
	if err != nil {
		return item, err
	}

	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	template := template.New(viewModel.Type)
	template.Parse(templateText)
	template.Execute(writer, viewModel)

	return item, nil
}
