// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"net/http"
)

func Error(headerWriter header.HeaderWriter, templateProvider templates.Provider, navigationOrchestrator *orchestrator.NavigationOrchestrator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)
		w.WriteHeader(http.StatusNotFound)

		// get the error template
		hostname := getBaseUrlFromRequest(r)
		errorTemplate, err := templateProvider.GetFullTemplate(hostname, templates.ErrorTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// create the view model
		errorModel := viewmodel.Model{
			Content: "",
		}

		errorModel.Type = "error"
		errorModel.Title = "Not found"
		errorModel.Description = "The requested resource was not found."
		errorModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		errorModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		// render the template
		renderTemplate(errorTemplate, errorModel, w)
	})
}
