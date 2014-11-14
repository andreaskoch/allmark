// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/services/thumbnail"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"io"
	"net/http"
)

type Item struct {
	logger logger.Logger

	fileOrchestrator      orchestrator.FileOrchestrator
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
	templateProvider      templates.Provider

	thumbnailIndex *thumbnail.Index

	error404Handler Handler
}

func (handler *Item) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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
			header.ContentType(w, r, "text/html")
			header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
			header.ETag(w, r, model.Hash)

			handler.render(w, model)
			return
		}

		// stage 2: check if there is a file for the request
		if file, found := handler.fileOrchestrator.GetFile(requestRoute); found {

			handler.logger.Info("Returning file %q", requestRoute)

			// set  headers
			header.ContentType(w, r, file.MimeType)
			header.Cache(w, r, header.STATICCONTENT_CACHEDURATION_SECONDS)
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

		// stage 3: thumbnails
		if baseRoute, dimensions := thumbnail.GetThumbnailDimensionsFromRoute(requestRoute); baseRoute.Value() != requestRoute.Value() {

			handler.logger.Debug("Requesting a thumbnail for route %q (Dimensions: %s)", baseRoute.Value(), dimensions)

			// get the thumbs for the base route
			if thumbs, thumbsExist := handler.thumbnailIndex.GetThumbs(baseRoute.Value()); thumbsExist {

				// get the thumb for the supplied dimensions
				if matchingThumb, matchingThumbExists := thumbs.GetThumbBySize(dimensions); matchingThumbExists {

					thumbnailFilePath := handler.thumbnailIndex.GetThumbnailFilepath(matchingThumb)
					http.ServeFile(w, r, thumbnailFilePath)

					return
				}
			}
		}

		handler.logger.Warn("No item or file found for route %q", requestRoute)

		// display a 404 error page
		error404Handler := handler.error404Handler.Func()
		error404Handler(w, r)
	}
}

func (handler *Item) render(writer io.Writer, viewModel viewmodel.Model) {

	// get a template
	template, err := handler.templateProvider.GetFullTemplate(viewModel.Type)
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
