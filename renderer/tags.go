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

	// relative file path provider
	relativePath := func(item *repository.Item) string {
		return item.AbsolutePath
	}

	// absolute file path provider
	absolutePath := func(item *repository.Item) string {
		return item.AbsolutePath
	}

	// content converter
	content := func(item *repository.Item) string {
		return ""
	}

	// render the tagmap content
	tagmapModel := mapper.MapTagmap(tags, tagPath, relativePath, absolutePath, content)
	tagmapContent := renderer.renderTagmap(tagmapContentTemplate, &tagmapModel)

	tagmapPageModel := view.Model{
		Title:                "Tags",
		Description:          "A list of all tags in this repository.",
		Content:              tagmapContent,
		ToplevelNavigation:   renderer.root.ToplevelNavigation,
		BreadcrumbNavigation: renderer.root.BreadcrumbNavigation,
		Type:                 "tagmap",
	}

	writeTemplate(tagmapPageModel, tagmapTemplate, writer)
}

func (renderer *Renderer) renderTagmap(templ *template.Template, tagmapModel *view.TagMap) string {

	// render
	buffer := new(bytes.Buffer)
	writeTemplate(tagmapModel, templ, buffer)

	// get the produced html code
	rootCode := buffer.String()

	return rootCode
}
