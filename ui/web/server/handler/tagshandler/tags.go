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

	templateProvider := templates.NewProvider(".")

	tagPathProvider := patherFactory.Absolute("/tags.html#")
	tagsOrchestrator := orchestrator.NewTagsOrchestrator(itemIndex, tagPathProvider)

	return &TagsHandler{
		logger:           logger,
		itemIndex:        itemIndex,
		config:           config,
		patherFactory:    patherFactory,
		templateProvider: templateProvider,
		tagsOrchestrator: tagsOrchestrator,
	}
}

type TagsHandler struct {
	logger           logger.Logger
	itemIndex        *index.ItemIndex
	config           *config.Config
	patherFactory    paths.PatherFactory
	templateProvider *templates.Provider
	tagsOrchestrator orchestrator.TagsOrchestrator
}

func (handler *TagsHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		tagmapTemplate, err := handler.templateProvider.GetFullTemplate(templates.TagmapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		tagmapContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.TagmapContentTemplateName)
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

		tagmapViewModel := viewmodel.Model{
			Content: tagMapItems,
		}

		tagmapViewModel.Title = "Tags"
		tagmapViewModel.Description = "A list of all tags in this repository."
		tagmapViewModel.ToplevelNavigation = orchestrator.GetToplevelNavigation(handler.itemIndex)
		tagmapViewModel.Type = "tagmap"

		handlerutil.RenderTemplate(tagmapViewModel, tagmapTemplate, w)
	}
}

func renderTagmapEntry(templ *template.Template, tagModel *viewmodel.Tag) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(tagModel, templ, buffer)
	return buffer.String()
}
