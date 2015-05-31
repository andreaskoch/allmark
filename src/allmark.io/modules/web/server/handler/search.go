// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	html "html/template"
	"net/http"
	"strings"
	"text/template"

	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
)

type Search struct {
	logger                 logger.Logger
	headerWriter           header.HeaderWriter
	navigationOrchestrator *orchestrator.NavigationOrchestrator
	searchOrchestrator     *orchestrator.SearchOrchestrator
	templateProvider       templates.Provider
	error404Handler        Handler
}

func (handler *Search) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseUrlFromRequest(r)

		// get the query parameter
		query, _ := getQueryParameterFromUrl(*r.URL)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the search template
		searchTemplate, err := handler.templateProvider.GetFullTemplate(hostname, templates.SearchTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page parameters
		pageType := "search"
		headline := getPageTitle(query)
		pageTitle := handler.searchOrchestrator.GetPageTitle(headline)
		description := getDescription(query)

		// Page model
		viewModel := viewmodel.Model{}
		viewModel.Type = pageType
		viewModel.Title = headline
		viewModel.PageTitle = pageTitle
		viewModel.Description = description
		viewModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		viewModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		// get the search result content template
		searchResultContentTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.SearchContentTemplateName)
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
