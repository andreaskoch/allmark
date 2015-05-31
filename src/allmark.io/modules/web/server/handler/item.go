// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"io"
	"net/http"
)

type Item struct {
	logger                logger.Logger
	headerWriter          header.HeaderWriter
	fileOrchestrator      *orchestrator.FileOrchestrator
	viewModelOrchestrator *orchestrator.ViewModelOrchestrator
	templateProvider      templates.Provider
	error404Handler       Handler
}

func (handler *Item) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		baseUrl := getBaseUrlFromRequest(r)

		// get the request route
		requestRoute := getRouteFromRequest(r)

		// make sure the request body is closed
		defer r.Body.Close()

		handler.logger.Info("Requesting %q", requestRoute)

		// stage 1: check if there is a item for the request
		if model, found := handler.viewModelOrchestrator.GetFullViewModel(requestRoute); found {

			handler.logger.Info("Returning item %q", requestRoute)

			// set headers
			handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)
			header.ETag(w, model.Hash)

			handler.render(w, baseUrl, model)
			return
		}

		// stage 2: check if there is a file for the request
		if file, found := handler.fileOrchestrator.GetFile(requestRoute); found {

			handler.logger.Info("Returning file %q", requestRoute)

			// set  headers
			handler.headerWriter.Write(w, file.MimeType)
			header.ETag(w, file.Hash)

			// get the content provider
			contentProvider := handler.fileOrchestrator.GetFileContentProvider(requestRoute)
			if contentProvider == nil {
				handler.logger.Error("There is no content provider for file %q", requestRoute)
			}

			filename := file.Name
			lastModifiedTime := file.LastModified

			contentProvider.Data(func(content io.ReadSeeker) error {
				http.ServeContent(w, r, filename, lastModifiedTime, content)
				return nil
			})

			return
		}

		handler.logger.Warn("No item or file found for route %q", requestRoute)

		// display a 404 error page
		error404Handler := handler.error404Handler.Func()
		error404Handler(w, r)
	}
}

func (handler *Item) render(writer io.Writer, baseUrl string, viewModel viewmodel.Model) {

	// get a template
	templateName := viewModel.Type
	template, err := handler.templateProvider.GetFullTemplate(baseUrl, templateName)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", templateName)
		return
	}

	// render template
	if err := renderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return
	}

}
