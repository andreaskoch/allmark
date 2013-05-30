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

	indexer          *repository.Indexer
	repositoryPath   string
	pathProvider     *path.Provider
	templateProvider *templates.Provider
	config           *config.Config
}

func New(repositoryPath string, config *config.Config, useTempDir bool) *Renderer {

	// create an index from the repository
	indexer, err := repository.New(repositoryPath, config, useTempDir)
	if err != nil {
		panic(fmt.Sprintf("Cannot create an item from folder %q.\nError: %s\n", repositoryPath, err))
	}

	return &Renderer{
		Rendered: make(chan *repository.Item),
		Removed:  make(chan *repository.Item),

		indexer:          indexer,
		repositoryPath:   repositoryPath,
		pathProvider:     path.NewProvider(repositoryPath, useTempDir),
		templateProvider: templates.NewProvider(config.TemplatesFolder()),
		config:           config,
	}

}

func (renderer *Renderer) Execute() {

	// start the indexer
	renderer.indexer.Execute()

	// render new items as they come in
	go func() {
		for {
			select {
			case item := <-renderer.indexer.New:

				// render the items
				fmt.Printf("Rendering item %q\n", item)
				renderer.renderItem(item)

				// attach change listeners
				renderer.listenForChanges(item)

			case item := <-renderer.indexer.Deleted:

				// remove the item
				fmt.Printf("Removing item %q\n", item)
				renderer.removeItem(item)
			}
		}
	}()

}

func (renderer *Renderer) listenForChanges(item *repository.Item) {
	go func() {
		for {
			select {
			case <-item.Modified:
				fmt.Printf("Rendering item %q\n", item)
				renderer.renderItem(item)

			case <-item.Moved:
				fmt.Printf("Removing item %q\n", item)
				renderer.removeItem(item)
			}
		}
	}()
}

func (renderer *Renderer) removeItem(item *repository.Item) {

	targetPath := renderer.pathProvider.GetRenderTargetPath(item)

	go func() {
		fmt.Printf("Removing %q\n", targetPath)
		os.Remove(targetPath)

		renderer.Removed <- item
	}()

}

func (renderer *Renderer) renderItem(item *repository.Item) {

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
