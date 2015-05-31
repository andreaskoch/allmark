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
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Print struct {
	logger                     logger.Logger
	headerWriter               header.HeaderWriter
	converterModelOrchestrator *orchestrator.ConversionModelOrchestrator
	templateProvider           templates.Provider
	error404Handler            Handler
}

func (handler *Print) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "print" or ".print" suffix from the path
		path = strings.TrimSuffix(path, "print")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// check if there is a item for the request
		baseUrl := getBaseUrlFromRequest(r)
		viewModel, found := handler.converterModelOrchestrator.GetConversionModel(baseUrl, requestRoute)
		if !found {

			// display a 404 error page
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		// render the view model
		handler.render(w, baseUrl, viewModel)
	}
}

func (handler *Print) render(writer io.Writer, baseUrl string, viewModel viewmodel.ConversionModel) {

	// get a template
	template, err := handler.templateProvider.GetSubTemplate(baseUrl, templates.ConversionTemplateName)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", viewModel.Type)
		return
	}

	// render template
	if err := renderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return
	}

}
