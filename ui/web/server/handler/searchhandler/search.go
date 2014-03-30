// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package searchhandler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
	"text/template"
)

var itemsPerPage = 5

func New(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex, fullTextIndex *search.FullTextIndex) *SearchHandler {

	// templates
	templateProvider := templates.NewProvider(".")

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	// search
	searchResultPathProvider := patherFactory.Absolute("/")
	searchOrchestrator := orchestrator.NewSearchOrchestrator(fullTextIndex, searchResultPathProvider)

	return &SearchHandler{
		logger:                 logger,
		itemIndex:              itemIndex,
		config:                 config,
		patherFactory:          patherFactory,
		templateProvider:       templateProvider,
		navigationOrchestrator: navigationOrchestrator,
		searchOrchestrator:     searchOrchestrator,
	}
}

type SearchHandler struct {
	logger                 logger.Logger
	itemIndex              *index.ItemIndex
	config                 *config.Config
	patherFactory          paths.PatherFactory
	templateProvider       *templates.Provider
	navigationOrchestrator orchestrator.NavigationOrchestrator
	searchOrchestrator     orchestrator.SearchOrchestrator
}

func (handler *SearchHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the query parameter
		query, _ := handlerutil.GetQueryParameterFromUrl(*r.URL)

		// read the page url-parameter
		page, pageParameterIsAvailable := handlerutil.GetPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the search template
		searchTemplate, err := handler.templateProvider.GetFullTemplate(templates.SearchTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// get the search result content template
		searchResultContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.SearchContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// root entry / channel item
		searchResultModel := handler.searchOrchestrator.GetSearchResults(query, page)
		searchResults := render(searchResultContentTemplate, searchResultModel)

		sitemapPageModel := viewmodel.Model{
			Content: searchResults,
		}

		sitemapPageModel.Type = "search"
		sitemapPageModel.Title = "Search"
		sitemapPageModel.Description = "Search results"
		sitemapPageModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		sitemapPageModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(handler.itemIndex.Root())

		handlerutil.RenderTemplate(sitemapPageModel, searchTemplate, w)
	}
}

func render(templ *template.Template, searchModel viewmodel.Search) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(searchModel, templ, buffer)
	return buffer.String()
}
