// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

func Print(logger logger.Logger,
	headerWriter header.HeaderWriter,
	conversionModelOrchestrator *orchestrator.ConversionModelOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	render := func(writer io.Writer, baseUrl string, viewModel viewmodel.ConversionModel) {

		// get a template
		template, err := templateProvider.GetSubTemplate(baseUrl, templates.ConversionTemplateName)
		if err != nil {
			logger.Error("No template for item of type %q.", viewModel.Type)
			return
		}

		// render template
		if err := renderTemplate(template, viewModel, writer); err != nil {
			logger.Error("%s", err)
			return
		}

	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)

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
		viewModel, found := conversionModelOrchestrator.GetConversionModel(baseUrl, requestRoute)
		if !found {

			// display a 404 error page
			error404Handler.ServeHTTP(w, r)
			return
		}

		// render the view model
		render(w, baseUrl, viewModel)
	})
}
