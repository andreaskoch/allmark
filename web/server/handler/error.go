// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"net/http"
)

type Error struct {
	logger logger.Logger

	templateProvider       templates.Provider
	navigationOrchestrator *orchestrator.NavigationOrchestrator
}

func (handler Error) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/html; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)

		// get the error template
		errorTemplate, err := handler.templateProvider.GetFullTemplate(templates.ErrorTemplateName)
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

		// set 404 status code
		w.WriteHeader(http.StatusNotFound)

		// render the template
		renderTemplate(errorModel, errorTemplate, w)
	}
}
