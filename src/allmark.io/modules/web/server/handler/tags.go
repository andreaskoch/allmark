// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
)

type Tags struct {
	logger                 logger.Logger
	headerWriter           header.HeaderWriter
	navigationOrchestrator *orchestrator.NavigationOrchestrator
	tagsOrchestrator       *orchestrator.TagsOrchestrator
	templateProvider       templates.Provider
}

func (handler *Tags) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseUrlFromRequest(r)

		tagmapTemplate, err := handler.templateProvider.GetFullTemplate(hostname, templates.TagmapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page conent
		tagmapContentTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.TagmapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		tagMapItems := ""
		tags := handler.tagsOrchestrator.GetTags()

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
		pageTitle := handler.tagsOrchestrator.GetPageTitle(headline)

		// Page model
		tagmapViewModel := viewmodel.Model{
			Content: tagMapItems,
		}

		tagmapViewModel.Type = pageType
		tagmapViewModel.Title = headline
		tagmapViewModel.PageTitle = pageTitle
		tagmapViewModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		tagmapViewModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(route.New())
		tagmapViewModel.TagCloud = handler.tagsOrchestrator.GetTagCloud()

		renderTemplate(tagmapViewModel, tagmapTemplate, w)
	}
}

func renderTagmapEntry(templ *template.Template, tagModel *viewmodel.Tag) string {
	buffer := new(bytes.Buffer)
	renderTemplate(tagModel, templ, buffer)
	return buffer.String()
}
