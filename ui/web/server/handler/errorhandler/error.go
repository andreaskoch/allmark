// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errorhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) *ErrorHandler {

	// templates
	templateProvider := templates.NewProvider(".")

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	return &ErrorHandler{
		logger:                 logger,
		itemIndex:              itemIndex,
		config:                 config,
		templateProvider:       templateProvider,
		navigationOrchestrator: navigationOrchestrator,
	}
}

type ErrorHandler struct {
	logger                 logger.Logger
	itemIndex              *index.ItemIndex
	config                 *config.Config
	templateProvider       *templates.Provider
	navigationOrchestrator orchestrator.NavigationOrchestrator
}

func (handler *ErrorHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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
		errorModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(handler.itemIndex.Root())

		// set 404 status code
		w.WriteHeader(http.StatusNotFound)

		// render the template
		handlerutil.RenderTemplate(errorModel, errorTemplate, w)
	}
}
