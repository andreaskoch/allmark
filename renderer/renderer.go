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
	"os"
	"text/template"
)

type Renderer struct {
	Rendered chan *repository.Item
	Removed  chan *repository.Item

	repositoryPath   string
	pathProvider     *path.Provider
	templateProvider *templates.Provider
	config           *config.Config
}

func New(repositoryPath string, config *config.Config, useTempDir bool) *Renderer {

	return &Renderer{
		Rendered: make(chan *repository.Item),
		Removed:  make(chan *repository.Item),

		repositoryPath:   repositoryPath,
		pathProvider:     path.NewProvider(repositoryPath, useTempDir),
		templateProvider: templates.NewProvider(config.TemplatesFolder()),
		config:           config,
	}

}

func (renderer *Renderer) Execute() {

	// create an index from the repository
	root, err := repository.NewRoot(renderer.repositoryPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot create an item from folder %q.\nError: %s\n", renderer.repositoryPath, err))
	}

	// start a change listener
	defer renderer.listenForChanges(root)

	// start rendering the root item
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

}

func (renderer *Renderer) listenForChanges(item *repository.Item) {

	for _, child := range item.Childs {
		renderer.listenForChanges(child) // recurse
	}

	go func() {
		for {
			select {
			case <-item.Modified:
				fmt.Printf("Rendering item %s\n", item)
				renderer.renderItem(item)

			case <-item.Moved:
				fmt.Printf("Item %s has moved", item)
				renderer.removeItem(item)
			}
		}
	}()

}

func (renderer *Renderer) removeItem(item *repository.Item) {

	// remove all childs first
	for _, child := range item.Childs {
		renderer.removeItem(child) // recurse
	}

	targetPath := renderer.pathProvider.GetRenderTargetPath(item)

	go func() {
		//os.Remove(targetPath)
		fmt.Printf("Removing %q\n", targetPath)

		renderer.Removed <- item
	}()

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

		// pass along
		go func() {
			renderer.Rendered <- item
		}()

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
