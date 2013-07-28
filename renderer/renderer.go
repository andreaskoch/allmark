// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"bytes"
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
	"sort"
	"strings"
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

				// Sort childs items
				renderer.root.Sort()

				// render
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

					// Sort childs items
					renderer.root.Sort()

					// render
					renderer.renderRecursive(renderer.root)
				}

			}
		}
	}()

}

func (renderer *Renderer) Error404(writer io.Writer) {

	// get the 404 page template
	templateType := templates.ErrorTemplateName
	template, err := renderer.templateProvider.GetFullTemplate(templateType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No template of type %s found.", templateType)
		return
	}

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
}

func (renderer *Renderer) Sitemap(writer io.Writer) {

	if renderer.root == nil {
		fmt.Println("The root is not ready yet.")
		return
	}

	// get the sitemap content template
	sitemapContentTemplate, err := renderer.templateProvider.GetSubTemplate(templates.SitemapContentTemplateName)
	if err != nil {
		return
	}

	// get the sitemap template
	sitemapTemplate, err := renderer.templateProvider.GetFullTemplate(templates.SitemapTemplateName)
	if err != nil {
		return
	}

	// render the sitemap content
	sitemapContentModel := mapper.MapSitemap(renderer.root)
	sitemapContent := renderer.renderSitemapEntry(sitemapContentTemplate, sitemapContentModel)

	sitemapPageModel := view.Model{
		Title:                "Sitemap",
		Description:          "A list of all items in this repository.",
		Content:              sitemapContent,
		ToplevelNavigation:   renderer.root.ToplevelNavigation,
		BreadcrumbNavigation: renderer.root.BreadcrumbNavigation,
		Type:                 "sitemap",
	}

	writeTemplate(sitemapPageModel, sitemapTemplate, writer)
}

func (renderer *Renderer) RSS(writer io.Writer) {

	fmt.Fprintln(writer, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintln(writer, `<rss version="2.0">`)
	fmt.Fprintln(writer, `<channel>`)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, fmt.Sprintf(`<title><![CDATA[%s]]></title>`, renderer.root.Title))
	fmt.Fprintln(writer, fmt.Sprintf(`<description><![CDATA[%s]]></description>`, renderer.root.Description))
	fmt.Fprintln(writer, fmt.Sprintf(`<link>%s</link>`, getItemLocation(renderer.root)))
	fmt.Fprintln(writer, fmt.Sprintf(`<pubData>%s</pubData>`, getItemDate(renderer.root)))
	fmt.Fprintln(writer)

	childsByDate := getAllItemsByDate(renderer.root)
	for _, item := range childsByDate {
		fmt.Fprintln(writer, `<item>`)
		fmt.Fprintln(writer, fmt.Sprintf(`<title><![CDATA[%s]]></title>`, item.Title))
		fmt.Fprintln(writer, fmt.Sprintf(`<description><![CDATA[%s]]></description>`, item.Description))
		fmt.Fprintln(writer, fmt.Sprintf(`<link>%s</link>`, getItemLocation(item)))
		fmt.Fprintln(writer, fmt.Sprintf(`<pubData>%s</pubData>`, getItemDate(item)))
		fmt.Fprintln(writer, `</item>`)
		fmt.Fprintln(writer)
	}

	fmt.Fprintln(writer, `</channel>`)
	fmt.Fprintln(writer, `</rss>`)

}

func getItemDate(item *repository.Item) string {
	return item.Date
}

func getItemLocation(item *repository.Item) string {
	route := item.AbsoluteRoute
	location := fmt.Sprintf(`http://%s/%s`, "example.com", route)
	return location
}

func getAllItemsByDate(root *repository.Item) repository.Items {
	childs := repository.GetAllChilds(root)
	sort.Sort(childs)
	return childs
}

func (renderer *Renderer) renderSitemapEntry(templ *template.Template, sitemapModel *view.Sitemap) string {

	// render
	buffer := new(bytes.Buffer)
	writeTemplate(sitemapModel, templ, buffer)

	// get the produced html code
	rootCode := buffer.String()

	if len(sitemapModel.Childs) > 0 {

		// render all childs

		childCode := ""
		for _, child := range sitemapModel.Childs {
			childCode += "\n" + renderer.renderSitemapEntry(templ, child)
		}

		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, childCode, 1)

	} else {

		// no childs
		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, "", 1)

	}

	return rootCode
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
	if template, err := renderer.templateProvider.GetFullTemplate(item.Type); err == nil {

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
	mapper.MapItem(item)
}

func writeTemplate(model interface{}, template *template.Template, writer io.Writer) {
	err := template.Execute(writer, model)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
