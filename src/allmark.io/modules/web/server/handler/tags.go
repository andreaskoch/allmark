// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

type Tags struct {
	logger logger.Logger

	navigationOrchestrator *orchestrator.NavigationOrchestrator
	tagsOrchestrator       *orchestrator.TagsOrchestrator

	templateProvider templates.Provider
}

func (self *Tags) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/html; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		hostname := getHostnameFromRequest(r)

		tagmapTemplate, err := self.templateProvider.GetFullTemplate(hostname, templates.TagmapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		tagmapContentTemplate, err := self.templateProvider.GetSubTemplate(hostname, templates.TagmapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		tagMapItems := ""
		tags := self.tagsOrchestrator.GetTags()

		if len(tags) > 0 {
			for _, tag := range tags {
				tagMapItems += renderTagmapEntry(tagmapContentTemplate, tag)
			}
		} else {
			tagMapItems = "-- There are currently not tagged items --"
		}

		tagmapViewModel := viewmodel.Model{
			Content: tagMapItems,
		}

		tagmapViewModel.Type = "tagmap"
		tagmapViewModel.Title = "Tags"
		tagmapViewModel.ToplevelNavigation = self.navigationOrchestrator.GetToplevelNavigation()
		tagmapViewModel.BreadcrumbNavigation = self.navigationOrchestrator.GetBreadcrumbNavigation(route.New())
		tagmapViewModel.TagCloud = self.tagsOrchestrator.GetTagCloud()

		renderTemplate(tagmapViewModel, tagmapTemplate, w)
	}
}

func renderTagmapEntry(templ *template.Template, tagModel *viewmodel.Tag) string {
	buffer := new(bytes.Buffer)
	renderTemplate(tagModel, templ, buffer)
	return buffer.String()
}
