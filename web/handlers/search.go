// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"fmt"
	html "html/template"
	"net/http"
	"strings"
	"text/template"

	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/orchestrator"
	"github.com/andreaskoch/allmark/web/view/templates"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
)

func Search(headerWriter header.HeaderWriter,
	navigationOrchestrator *orchestrator.NavigationOrchestrator,
	searchOrchestrator *orchestrator.SearchOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseURLFromRequest(r)

		// get the query parameter
		query, _ := getQueryParameterFromURL(*r.URL)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromURL(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the search template
		searchTemplate, err := templateProvider.GetSearchTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page parameters
		pageType := "search"
		headline := getPageTitle(query)
		pageTitle := searchOrchestrator.GetPageTitle(headline)
		description := getDescription(query)

		// Page model
		pageModel := viewmodel.Model{}
		pageModel.Type = pageType
		pageModel.Title = headline
		pageModel.PageTitle = pageTitle
		pageModel.Description = description
		pageModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		pageModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		// get the search results
		searchResultsModel := searchOrchestrator.GetSearchResults(query, page)

		// display error 404 non-existing page has been requested
		if searchResultsModel.ResultCount == 0 && page > 1 {
			error404Handler.ServeHTTP(w, r)
			return
		}

		// assemble the page model
		searchResultPage := viewmodel.Search{}
		searchResultPage.Model = pageModel
		searchResultPage.Results = searchResultsModel

		renderTemplate(searchTemplate, searchResultPage, w)
	})

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

func renderSearchResultModel(templ *template.Template, searchModel viewmodel.Search) string {
	buffer := new(bytes.Buffer)
	renderTemplate(templ, searchModel, buffer)
	return buffer.String()
}
