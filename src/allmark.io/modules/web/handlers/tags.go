// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"allmark.io/modules/common/route"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
)

func Tags(headerWriter header.HeaderWriter,
	navigationOrchestrator *orchestrator.NavigationOrchestrator,
	tagsOrchestrator *orchestrator.TagsOrchestrator,
	templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseURLFromRequest(r)

		tagmapTemplate, err := templateProvider.GetFullTemplate(hostname, templates.TagmapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page conent
		tagmapContentTemplate, err := templateProvider.GetSubTemplate(hostname, templates.TagmapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		tagMapItems := ""
		tags := tagsOrchestrator.GetTags()

		if len(tags) > 0 {
			for _, tag := range tags {
				tagMapItems += renderTagmapEntry(tagmapContentTemplate, tag)
			}
		} else {
			tagMapItems = "-- There are currently not tagged items --"
		}

		// Page parameters
		pageType := "tagmap"
		headline := "Tags"
		pageTitle := tagsOrchestrator.GetPageTitle(headline)

		// Page model
		tagmapViewModel := viewmodel.Model{
			Content: tagMapItems,
		}

		tagmapViewModel.Type = pageType
		tagmapViewModel.Title = headline
		tagmapViewModel.PageTitle = pageTitle
		tagmapViewModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		tagmapViewModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())
		tagmapViewModel.TagCloud = tagsOrchestrator.GetTagCloud()

		renderTemplate(tagmapTemplate, tagmapViewModel, w)
	})
}

func renderTagmapEntry(templ *template.Template, tagModel *viewmodel.Tag) string {
	buffer := new(bytes.Buffer)
	renderTemplate(templ, tagModel, buffer)
	return buffer.String()
}
