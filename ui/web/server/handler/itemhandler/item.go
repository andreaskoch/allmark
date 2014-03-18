// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) *ItemHandler {

	templateProvider := templates.NewProvider(".")
	viewModelOrchestrator := orchestrator.NewViewModelOrchestrator(itemIndex, converter)
	error404Handler := errorhandler.New(logger, config, itemIndex)

	rewrites := []RequestRewrite{
		NewRewrite("^favicon.ico", "theme/favicon.ico"),
	}

	return &ItemHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		fileIndex:             fileIndex,
		config:                config,
		patherFactory:         patherFactory,
		templateProvider:      templateProvider,
		error404Handler:       error404Handler,
		viewModelOrchestrator: viewModelOrchestrator,
		rewrites:              rewrites,
	}
}

type ItemHandler struct {
	logger                logger.Logger
	itemIndex             *index.ItemIndex
	fileIndex             *index.FileIndex
	config                *config.Config
	patherFactory         paths.PatherFactory
	templateProvider      *templates.Provider
	error404Handler       *errorhandler.ErrorHandler
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
	rewrites              []RequestRewrite
}

func (handler *ItemHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestRoute, err := handlerutil.GetRouteFromRequest(r)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check for theme files
		themeRoute, err := route.NewFromRequest("theme")
		if err != nil {
			fmt.Fprintln(w, "%s", err)
			return
		}

		if isThemeFile := requestRoute.IsChildOf(themeRoute); isThemeFile {

			if file, found := handler.fileIndex.IsMatch(*requestRoute); found {
				handler.serveContent(requestRoute.Value(), file.ContentProvider(), w)
				return
			}
		}

		// stage 2: check if there is a item for the request
		if item, found := handler.itemIndex.IsMatch(*requestRoute); found {

			// create the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := handler.viewModelOrchestrator.GetViewModel(pathProvider, item)

			// render the view model
			handler.render(w, viewModel)
			return
		}

		// stage 3: check if there is a file for the request
		if file, found := handler.itemIndex.IsFileMatch(*requestRoute); found {
			handler.serveContent(requestRoute.Value(), file.ContentProvider(), w)
			return
		}

		// stage 4: check if there is rewrite for the the request
		for _, rewrite := range handler.rewrites {

			isMatch, target := rewrite.Match(*requestRoute)
			if !isMatch {
				continue
			}

			if file, found := handler.fileIndex.IsMatch(target); found {
				handler.serveContent(target.Value(), file.ContentProvider(), w)
				return
			}
		}

		// display a 404 error page
		error404Handler := handler.error404Handler.Func()
		error404Handler(w, r)
	}
}

func (handler *ItemHandler) render(writer io.Writer, viewModel viewmodel.Model) {

	// get a template
	template, err := handler.templateProvider.GetFullTemplate(viewModel.Type)
	if err != nil {
		handler.logger.Error("No template for item of type %q.", viewModel.Type)
		return
	}

	// render template
	if err := handlerutil.RenderTemplate(viewModel, template, writer); err != nil {
		handler.logger.Error("%s", err)
		return
	}

}

func (handler *ItemHandler) serveContent(filename string, contentProvider *content.ContentProvider, w http.ResponseWriter) {
	// read the file data
	data, err := contentProvider.Data()
	if err != nil {
		handler.logger.Error("Unable to read the content: %s", err)
		return
	}

	// content type detection
	// derive content type from file extension
	fileExtension := filepath.Ext(filename)
	contentType := mime.TypeByExtension(fileExtension)
	if contentType == "" {
		// fallback: derive content type from data
		contentType = http.DetectContentType(data)
	}

	// set headers
	w.Header().Set("Content-Type", contentType)

	fmt.Fprintf(w, "%s", data)
}
