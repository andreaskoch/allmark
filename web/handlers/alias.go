// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"github.com/elWyatt/allmark/web/view/templates"
	"github.com/elWyatt/allmark/web/view/viewmodel"
	"fmt"
	"net/http"
	"strings"
)

// AliasLookup creates a http handler which redirects aliases to their documents.
func AliasLookup(
	headerWriter header.HeaderWriter,
	viewModelOrchestrator *orchestrator.ViewModelOrchestrator,
	fallbackHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// strip the "/!" prefix from the path
		path := r.URL.Path
		alias := strings.TrimPrefix(path, "/!")

		// locate the correct viewmodel for the given alias
		viewModel, found := viewModelOrchestrator.GetViewModelByAlias(alias)
		if !found {
			// no model found for the alias -> use fallback handler
			fallbackHandler.ServeHTTP(w, r)
			return
		}

		// determine the redirect url
		baseURL := getBaseURLFromRequest(r)
		redirectURL := baseURL + viewModel.BaseURL
		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	})

}

// AliasIndex creates a http handler which displays an index of all aliases.
func AliasIndex(
	headerWriter header.HeaderWriter,
	navigationOrchestrator *orchestrator.NavigationOrchestrator,
	aliasIndexOrchestrator *orchestrator.AliasIndexOrchestrator,
	templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseURLFromRequest(r)

		aliasIndexTemplate, err := templateProvider.GetAliasIndexTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// assemble the base view model
		title := "Shortlinks"
		description := "A list of all short links to different items in this repository."
		viewModel := viewmodel.Model{}

		viewModel.Type = "aliasindex"
		viewModel.Title = title
		viewModel.Description = description
		viewModel.PageTitle = aliasIndexOrchestrator.GetPageTitle(title)
		viewModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		viewModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		// assemble the specialized alias index viewmodel
		aliasIndexViewModel := viewmodel.AliasIndex{}
		aliasIndexViewModel.Model = viewModel
		aliasIndexViewModel.Aliases = aliasIndexOrchestrator.GetIndexEntries(hostname, "!")

		renderTemplate(aliasIndexTemplate, aliasIndexViewModel, w)

	})

}

type aliasIndexViewModel struct {
	viewmodel.Model
	Entries []viewmodel.Alias
}
