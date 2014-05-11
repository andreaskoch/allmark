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
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	html "html/template"
	"net/http"
	"strings"
	"text/template"
)

func New(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.Index, searcher *search.ItemSearch) *SearchHandler {

	// templates
	templateProvider := templates.NewProvider(config.TemplatesFolder())

	// errors
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	// search
	searchResultPathProvider := patherFactory.Absolute("/")
	searchOrchestrator := orchestrator.NewSearchOrchestrator(searcher, searchResultPathProvider)

	return &SearchHandler{
		logger:                 logger,
		itemIndex:              itemIndex,
		config:                 config,
		patherFactory:          patherFactory,
		templateProvider:       templateProvider,
		error404Handler:        error404Handler,
		navigationOrchestrator: navigationOrchestrator,
		searchOrchestrator:     searchOrchestrator,
	}
}

type SearchHandler struct {
	logger                 logger.Logger
	itemIndex              *index.Index
	config                 *config.Config
	patherFactory          paths.PatherFactory
	templateProvider       *templates.Provider
	error404Handler        *errorhandler.ErrorHandler
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

		viewModel := viewmodel.Model{}
		viewModel.Type = "search"
		viewModel.Title = getPageTitle(query)
		viewModel.Description = getDescription(query)
		viewModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		viewModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(handler.itemIndex.Root())

		// get the search result content template
		searchResultContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.SearchContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// get the search results
		searchResultModel := handler.searchOrchestrator.GetSearchResults(query, page)

		// display error 404 non-existing page has been requested
		if searchResultModel.ResultCount == 0 && page > 1 {
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		viewModel.Content = render(searchResultContentTemplate, searchResultModel)

		handlerutil.RenderTemplate(viewModel, searchTemplate, w)
	}
}

func getPageTitle(query string) string {
	if strings.TrimSpace(query) == "" {
		return "Search"
	}

	return fmt.Sprintf("%s | Search", html.HTMLEscapeString(query))
}

func getDescription(query string) string {
	if strings.TrimSpace(query) == "" {
		return "Search this repository."
	}

	return fmt.Sprintf("Search results for %q.", html.HTMLEscapeString(query))
}

func render(templ *template.Template, searchModel viewmodel.Search) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(searchModel, templ, buffer)
	return buffer.String()
}
