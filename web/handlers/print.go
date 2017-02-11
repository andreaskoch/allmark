// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/orchestrator"
	"github.com/andreaskoch/allmark/web/view/templates"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"io"
	"net/http"
	"strings"
)

func Print(logger logger.Logger,
	headerWriter header.HeaderWriter,
	conversionModelOrchestrator *orchestrator.ConversionModelOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	render := func(writer io.Writer, baseURL string, viewModel viewmodel.ConversionModel) {

		// get a template
		template, err := templateProvider.GetConversionTemplate(baseURL)
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

		// strip the "print" or ".print" suffix from the path
		path := r.URL.Path
		path = strings.TrimSuffix(path, "print")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// check if there is a item for the request
		baseURL := getBaseURLFromRequest(r)
		viewModel, found := conversionModelOrchestrator.GetConversionModel(baseURL, requestRoute)
		if !found {

			// display a 404 error page
			error404Handler.ServeHTTP(w, r)
			return
		}

		// render the view model
		render(w, baseURL, viewModel)
	})
}
