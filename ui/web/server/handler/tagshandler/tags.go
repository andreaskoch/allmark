// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tagshandler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
	"text/template"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) *TagsHandler {

	// templates
	templateProvider := templates.NewProvider(".")

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	// tags
	tagPathProvider := patherFactory.Absolute("/tags.html#")
	tagsOrchestrator := orchestrator.NewTagsOrchestrator(itemIndex, tagPathProvider)

	return &TagsHandler{
		logger:                 logger,
		itemIndex:              itemIndex,
		config:                 config,
		patherFactory:          patherFactory,
		templateProvider:       templateProvider,
		navigationOrchestrator: navigationOrchestrator,
		tagsOrchestrator:       tagsOrchestrator,
	}
}

type TagsHandler struct {
	logger                 logger.Logger
	itemIndex              *index.ItemIndex
	config                 *config.Config
	patherFactory          paths.PatherFactory
	templateProvider       *templates.Provider
	navigationOrchestrator orchestrator.NavigationOrchestrator
	tagsOrchestrator       orchestrator.TagsOrchestrator
}

func (self *TagsHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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
		tagmapViewModel.Description = "A list of all tags in this repository."
		tagmapViewModel.ToplevelNavigation = self.navigationOrchestrator.GetToplevelNavigation()
		tagmapViewModel.BreadcrumbNavigation = self.navigationOrchestrator.GetBreadcrumbNavigation(self.itemIndex.Root())
		tagmapViewModel.TagCloud = self.tagsOrchestrator.GetTagCloud()

		handlerutil.RenderTemplate(tagmapViewModel, tagmapTemplate, w)
	}
}

func renderTagmapEntry(templ *template.Template, tagModel *viewmodel.Tag) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(tagModel, templ, buffer)
	return buffer.String()
}
