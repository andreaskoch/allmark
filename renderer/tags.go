// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/templates"
	"github.com/andreaskoch/allmark/view"
	"io"
	"strings"
	"text/template"
)

func (renderer *Renderer) Tags(writer io.Writer, host string) {

	if renderer.root == nil {
		fmt.Fprintf(writer, "The root is not ready yet.")
		return
	}

	// get the tagmap content template
	tagmapContentTemplate, err := renderer.templateProvider.GetSubTemplate(templates.TagmapContentTemplateName)
	if err != nil {
		fmt.Fprintf(writer, "Content template not found. Error: %s", err)
		return
	}

	// get the tagmap template
	tagmapTemplate, err := renderer.templateProvider.GetFullTemplate(templates.TagmapTemplateName)
	if err != nil {
		fmt.Fprintf(writer, "Template not found. Error: %s", err)
		return
	}

	// render the tagmap content
	tagmapContentModel := mapper.MapSitemap(renderer.root)
	tagmapContent := renderer.renderTagmapEntry(tagmapContentTemplate, tagmapContentModel)

	sitemapPageModel := view.Model{
		Title:                "Tags",
		Description:          "A list of all tags in this repository.",
		Content:              tagmapContent,
		ToplevelNavigation:   renderer.root.ToplevelNavigation,
		BreadcrumbNavigation: renderer.root.BreadcrumbNavigation,
		Type:                 "tagmap",
	}

	writeTemplate(sitemapPageModel, tagmapTemplate, writer)
}

func (renderer *Renderer) renderTagmapEntry(templ *template.Template, sitemapModel *view.Sitemap) string {

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
