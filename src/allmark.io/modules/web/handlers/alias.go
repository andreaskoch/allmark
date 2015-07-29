// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

// AliasLookup creates a http handler which redirects aliases to their documents.
func AliasLookup(
	headerWriter header.HeaderWriter,
	viewModelOrchestrator *orchestrator.ViewModelOrchestrator,
	fallbackHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		alias := vars["alias"]

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

		aliasIndexTemplate, err := templateProvider.GetFullTemplate(hostname, templates.AliasIndexTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page conent
		aliasIndexContentTemplate, err := templateProvider.GetSubTemplate(hostname, templates.AliasIndexContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		aliasIndexContent := ""
		aliasIndexEntries := aliasIndexOrchestrator.GetIndexEntries(hostname, "!")

		if len(aliasIndexEntries) > 0 {
			for _, aliasIndexEntry := range aliasIndexEntries {
				aliasIndexContent += renderAliasIndexEntry(aliasIndexContentTemplate, aliasIndexEntry)
			}
		} else {
			aliasIndexContent = "-- There are currently not items with aliases in this repository --"
		}

		// Page model
		aliasIndexViewModel := viewmodel.Model{
			Content: aliasIndexContent,
		}

		title := "Shortlinks"
		description := "A list of all short links to different items in this repository."

		aliasIndexViewModel.Type = "aliasindex"
		aliasIndexViewModel.Title = title
		aliasIndexViewModel.Description = description
		aliasIndexViewModel.PageTitle = aliasIndexOrchestrator.GetPageTitle(title)
		aliasIndexViewModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		aliasIndexViewModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		renderTemplate(aliasIndexTemplate, aliasIndexViewModel, w)

	})

}

func renderAliasIndexEntry(templ *template.Template, model interface{}) string {
	buffer := new(bytes.Buffer)
	renderTemplate(templ, model, buffer)
	return buffer.String()
}
