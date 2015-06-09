// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"io"
	"net/http"
)

func Item(logger logger.Logger,
	headerWriter header.HeaderWriter,
	fileOrchestrator *orchestrator.FileOrchestrator,
	viewModelOrchestrator *orchestrator.ViewModelOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	render := func(writer io.Writer, baseUrl string, viewModel viewmodel.Model) {

		// get a template
		templateName := viewModel.Type
		template, err := templateProvider.GetFullTemplate(baseUrl, templateName)
		if err != nil {
			logger.Error("No template for item of type %q.", templateName)
			return
		}

		// render template
		if err := renderTemplate(template, viewModel, writer); err != nil {
			logger.Error("%s", err)
			return
		}

	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		baseUrl := getBaseUrlFromRequest(r)

		// get the request route
		requestRoute := getRouteFromRequest(r)

		// make sure the request body is closed
		defer r.Body.Close()

		logger.Info("Requesting %q", requestRoute)

		// stage 1: check if there is a item for the request
		if model, found := viewModelOrchestrator.GetFullViewModel(requestRoute); found {

			logger.Info("Returning item %q", requestRoute)

			// set headers
			headerWriter.Write(w, header.CONTENTTYPE_HTML)
			header.ETag(w, model.Hash)

			render(w, baseUrl, model)
			return
		}

		// stage 2: check if there is a file for the request
		if file, found := fileOrchestrator.GetFile(requestRoute); found {

			logger.Info("Returning file %q", requestRoute)

			// set  headers
			headerWriter.Write(w, file.MimeType)
			header.ETag(w, file.Hash)

			// get the content provider
			contentProvider := fileOrchestrator.GetFileContentProvider(requestRoute)
			if contentProvider == nil {
				logger.Error("There is no content provider for file %q", requestRoute)
			}

			filename := file.Name
			lastModifiedTime := file.LastModified

			contentProvider.Data(func(content io.ReadSeeker) error {
				http.ServeContent(w, r, filename, lastModifiedTime, content)
				return nil
			})

			return
		}

		logger.Warn("No item or file found for route %q", requestRoute)

		// display a 404 error page
		error404Handler.ServeHTTP(w, r)
	})
}
