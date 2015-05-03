// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"net/http"
)

type Error struct {
	logger                 logger.Logger
	headerWriter           header.HeaderWriter
	templateProvider       templates.Provider
	navigationOrchestrator *orchestrator.NavigationOrchestrator
}

func (handler Error) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)
		w.WriteHeader(http.StatusNotFound)

		// get the error template
		hostname := getHostnameFromRequest(r)
		errorTemplate, err := handler.templateProvider.GetFullTemplate(hostname, templates.ErrorTemplateName)
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
		errorModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		errorModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		// render the template
		renderTemplate(errorModel, errorTemplate, w)
	}
}
