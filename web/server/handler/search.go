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
	html "html/template"
	"net/http"
	"strings"
	"text/template"
)

type Search struct {
	logger logger.Logger

	navigationOrchestrator *orchestrator.NavigationOrchestrator
	searchOrchestrator     *orchestrator.SearchOrchestrator

	templateProvider templates.Provider

	error404Handler Handler
}

func (handler *Search) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/html; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		// get the query parameter
		query, _ := getQueryParameterFromUrl(*r.URL)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromUrl(*r.URL)
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
		viewModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(route.New())

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

		renderTemplate(viewModel, searchTemplate, w)
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
	renderTemplate(searchModel, templ, buffer)
	return buffer.String()
}
