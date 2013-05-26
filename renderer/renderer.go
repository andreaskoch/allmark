// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/watcher"
	"os"
	"text/template"
)

type Renderer struct {
	repositoryPath   string
	pathProvider     *path.Provider
	templateProvider *templates.Provider
	config           *config.Config
}

func New(repositoryPath string, config *config.Config, useTempDir bool) *Renderer {

	return &Renderer{
		repositoryPath:   repositoryPath,
		pathProvider:     path.NewProvider(repositoryPath, useTempDir),
		templateProvider: templates.NewProvider(config.TemplatesFolder()),
		config:           config,
	}

}

func (renderer *Renderer) Execute() *repository.ItemIndex {

	// create an index from the repository
	index, err := repository.NewItemIndex(renderer.repositoryPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot create an item index for folder %q.\nError: %s\n", renderer.repositoryPath, err))
	}

	root := index.Root()
	defer renderer.attachChangeListener(root)

	renderer.renderItem(root)

	// re-render on template change
	go func() {
		for {
			select {
			case event := <-renderer.templateProvider.TemplateChanged:
				fmt.Printf("Template %q changed. Rendering all items.\n", event.Filepath)
				renderer.renderItem(root)
			}
		}
	}()

	return index
}

func (renderer *Renderer) attachChangeListener(item *repository.Item) {

	for _, child := range item.Childs {

		// aggregate child events
		child.OnChange("Throw Item Events on Child Item change", func(event *watcher.WatchEvent) {
			item.Throw(event)
		})

		// recurse
		renderer.attachChangeListener(child)
	}

	// attach change listener
	item.OnChange("Render item on change", func(event *watcher.WatchEvent) {
		renderer.renderItem(item)
	})
}

func (renderer *Renderer) renderItem(item *repository.Item) {

	// render childs first
	for _, child := range item.Childs {
		renderer.renderItem(child) // recurse
	}

	// create the viewmodel
	mapper.Map(item, renderer.pathProvider)

	// get a template
	if template, err := renderer.templateProvider.GetTemplate(item.Type); err == nil {

		// render the template
		targetPath := renderer.pathProvider.GetRenderTargetPath(item)
		renderer.writeOutput(item, template, targetPath)

	} else {
		fmt.Fprintf(os.Stderr, "No template for item of type %q.", item.Type)
	}

}

func (renderer *Renderer) writeOutput(item *repository.Item, template *template.Template, targetPath string) {
	file, err := os.Create(targetPath)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	template.Execute(writer, item)
}
