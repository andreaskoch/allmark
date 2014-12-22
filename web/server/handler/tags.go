// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
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

		tagmapTemplate, err := self.templateProvider.GetFullTemplate(templates.TagmapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		tagmapContentTemplate, err := self.templateProvider.GetSubTemplate(templates.TagmapContentTemplateName)
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
