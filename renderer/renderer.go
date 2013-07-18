// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"io"
	"os"
	"text/template"
)

type Renderer struct {
	Rendered chan *repository.Item
	Removed  chan *repository.Item

	root *repository.Item

	rootIsReady      bool
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

		rootIsReady:      false,
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
				if !renderer.rootIsReady {
					fmt.Printf("Preparing %q\n", item)
					prepare(item)
				} else {
					fmt.Printf("Rendering %q\n", item)
					renderer.render(item)
				}

				// attach change listeners
				renderer.listenForChanges(item)

			case item := <-renderer.indexer.Deleted:

				// remove the item
				fmt.Printf("Removing %q\n", item)
				renderer.removeItem(item)
			}
		}
	}()

	go func() {
		for {
			select {
			case root := <-renderer.indexer.RootIsReady:

				// save the root item
				renderer.root = root

				// set the root is ready flag
				renderer.rootIsReady = true

				// render all items from the top
				fmt.Println("Root is ready. Rendering all items.")
				renderer.renderRecursive(root)
			}
		}
	}()

	// re-render on template change
	go func() {
		for {
			select {
			case <-renderer.templateProvider.Modified:

				if renderer.root != nil {
					fmt.Println("A template changed. Rendering all items.")
					renderer.renderRecursive(renderer.root)
				}

			}
		}
	}()

}

func (renderer *Renderer) Error404(writer io.Writer) {
	templateType := templates.ErrorTemplateName
	if template, err := renderer.templateProvider.GetTemplate(templateType); err == nil {

		// create a error view model
		title := "Not found"
		content := fmt.Sprintf("The requested item was not found.")
		errorModel := view.Error(title, content, renderer.root.RelativePath, renderer.root.AbsolutePath)

		// attach the toplevel navigation
		errorModel.ToplevelNavigation = renderer.root.ToplevelNavigation

		// attach the bread crumb navigation
		errorModel.BreadcrumbNavigation = renderer.root.BreadcrumbNavigation

		// render the template
		writeTemplate(errorModel, template, writer)

	} else {

		fmt.Fprintf(os.Stderr, "No template of type %s found.", templateType)

	}
}

func (renderer *Renderer) listenForChanges(item *repository.Item) {
	go func() {
		for {
			select {
			case <-item.Modified:
				fmt.Printf("Rendering %q\n", item)
				renderer.render(item)

				if parent := item.Parent; parent != nil {
					fmt.Printf("Rendering parent %q\n", parent)
					renderer.render(parent)
				}

			case <-item.Moved:
				fmt.Printf("Removing %q\n", item)
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

func (renderer *Renderer) renderRecursive(item *repository.Item) {
	for _, child := range item.Childs {
		renderer.renderRecursive(child)
	}

	renderer.render(item)
}

func (renderer *Renderer) render(item *repository.Item) {

	// prepare the item
	prepare(item)

	// render the bread crumb navigation
	attachBreadcrumbNavigation(item)

	// render the top-level navigation
	attachToplevelNavigation(renderer.root, item)

	// get a template
	if template, err := renderer.templateProvider.GetTemplate(item.Type); err == nil {

		// open the target file
		targetPath := renderer.pathProvider.GetRenderTargetPath(item)
		file, err := os.Create(targetPath)
		if err != nil {
			fmt.Errorf("%s", err)
		}

		writer := bufio.NewWriter(file)

		defer func() {
			writer.Flush()
		}()

		// render the template
		writeTemplate(item.Model, template, writer)

		// pass along
		go func() {
			renderer.Rendered <- item
		}()

	} else {

		fmt.Fprintf(os.Stderr, "No template for item of type %q.", item.Type)

	}

}

func prepare(item *repository.Item) {
	// parse the item
	parser.Parse(item)

	// convert the item
	converter.Convert(item)

	// create the viewmodel
	mapper.Map(item)
}

func writeTemplate(model *view.Model, template *template.Template, writer io.Writer) {
	template.Execute(writer, model)
}
