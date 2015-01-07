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
	logger logger.Logger

	fileOrchestrator      *orchestrator.FileOrchestrator
	viewModelOrchestrator *orchestrator.ViewModelOrchestrator
	templateProvider      templates.Provider

	error404Handler Handler
}

func (handler *Item) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		hostname := getHostnameFromRequest(r)

		// get the request route
		requestRoute, err := getRouteFromRequest(r)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		handler.logger.Info("Requesting %q", requestRoute)

		// stage 1: check if there is a item for the request
		if model, found := handler.viewModelOrchestrator.GetFullViewModel(requestRoute); found {

			handler.logger.Info("Returning item %q", requestRoute)

			// set headers
			header.ContentType(w, r, "text/html; charset=utf-8")
			header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
			header.VaryAcceptEncoding(w, r)
			header.ETag(w, r, model.Hash)

			handler.render(w, hostname, model)
			return
		}

		// stage 2: check if there is a file for the request
		if file, found := handler.fileOrchestrator.GetFile(requestRoute); found {

			handler.logger.Info("Returning file %q", requestRoute)

			// set  headers
			header.ContentType(w, r, file.MimeType)
			header.Cache(w, r, header.STATICCONTENT_CACHEDURATION_SECONDS)
			header.VaryAcceptEncoding(w, r)
			header.ETag(w, r, file.Hash)

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

func (handler *Item) render(writer io.Writer, hostname string, viewModel viewmodel.Model) {

	// get a template
	templateName := viewModel.Type
	template, err := handler.templateProvider.GetFullTemplate(hostname, templateName)
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
