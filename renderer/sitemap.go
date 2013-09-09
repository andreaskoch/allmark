// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"io"
	"strings"
	"text/template"
)

func (renderer *Renderer) Sitemap(writer io.Writer, host string) {

	targetFile := "sitemap.html"
	pathProvider := renderer.pathProvider
	rssRenderer := func(writer io.Writer, host string) {
		sitemap(writer, renderer.root, renderer.templateProvider)
	}

	cacheReponse(targetFile, pathProvider, rssRenderer, host, writer)
}

func sitemap(writer io.Writer, rootItem *repository.Item, templateProvider *templates.Provider) {

	if rootItem == nil {
		fmt.Fprintf(writer, "The root is not ready yet.")
		return
	}

	// get the sitemap content template
	sitemapContentTemplate, err := templateProvider.GetSubTemplate(templates.SitemapContentTemplateName)
	if err != nil {
		fmt.Fprintf(writer, "Content template not found. Error: %s", err)
		return
	}

	// get the sitemap template
	sitemapTemplate, err := templateProvider.GetFullTemplate(templates.SitemapTemplateName)
	if err != nil {
		fmt.Fprintf(writer, "Template not found. Error: %s", err)
		return
	}

	// render the sitemap content
	sitemapContentModel := mapper.MapSitemap(rootItem)
	sitemapContent := renderSitemapEntry(sitemapContentTemplate, sitemapContentModel)

	sitemapPageModel := view.Model{
		Title:                "Sitemap",
		Description:          "A list of all items in this repository.",
		Content:              sitemapContent,
		ToplevelNavigation:   rootItem.ToplevelNavigation,
		BreadcrumbNavigation: rootItem.BreadcrumbNavigation,
		Type:                 "sitemap",
	}

	writeTemplate(sitemapPageModel, sitemapTemplate, writer)
}

func renderSitemapEntry(templ *template.Template, sitemapModel *view.Sitemap) string {

	// render
	buffer := new(bytes.Buffer)
	writeTemplate(sitemapModel, templ, buffer)

	// get the produced html code
	rootCode := buffer.String()

	if len(sitemapModel.Childs) > 0 {

		// render all childs

		childCode := ""
		for _, child := range sitemapModel.Childs {
			childCode += "\n" + renderSitemapEntry(templ, child)
		}

		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, childCode, 1)

	} else {

		// no childs
		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, "", 1)

	}

	return rootCode
}
