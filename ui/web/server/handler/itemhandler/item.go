// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/updatehandler"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter, updateHandler *updatehandler.UpdateHandler) *ItemHandler {

	// templates
	templateProvider := templates.NewProvider(config.TemplatesFolder())

	// tags
	itemPathProvider := patherFactory.Absolute("/")
	tagPathProvider := patherFactory.Absolute("/tags.html#")
	tagsOrchestrator := orchestrator.NewTagsOrchestrator(itemIndex, tagPathProvider, itemPathProvider)

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// viewmodel
	viewModelOrchestrator := orchestrator.NewViewModelOrchestrator(itemIndex, converter, &navigationOrchestrator, &tagsOrchestrator)

	return &ItemHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		config:                config,
		patherFactory:         patherFactory,
		templateProvider:      templateProvider,
		error404Handler:       error404Handler,
		updateHandler:         updateHandler,
		viewModelOrchestrator: viewModelOrchestrator,
	}
}

type ItemHandler struct {
	logger                logger.Logger
	itemIndex             *index.Index
	config                *config.Config
	patherFactory         paths.PatherFactory
	templateProvider      *templates.Provider
	error404Handler       *errorhandler.ErrorHandler
	updateHandler         *updatehandler.UpdateHandler
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
}

func (handler *ItemHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	if handler.updateHandler != nil {

		handler.itemIndex.OnUpdate(func(item *model.Item) {

			handler.logger.Debug("Executing onUpdate handler for item %q", item)

			if item == nil {
				handler.logger.Error("The item is not initialzed.")
				return
			}

			pathProvider := handler.patherFactory.Absolute("/")
			viewModel := handler.viewModelOrchestrator.GetViewModel(pathProvider, item)
			updatehandler.UpdateMessage(viewModel)
		})

	}

	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestRoute, err := handlerutil.GetRouteFromRequest(r)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if item, found := handler.itemIndex.IsMatch(*requestRoute); found {

			// render the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := handler.viewModelOrchestrator.GetViewModel(pathProvider, item)
			handler.render(w, viewModel)
			return
		}

		// stage 2: check if there is a file for the request
		if file, found := handler.itemIndex.IsFileMatch(*requestRoute); found {

			filename := requestRoute.Value()

			lastModifiedTime, lastModErr := file.LastModified()
			if lastModErr != nil {
				handler.logger.Warn("Unable to determine lastmodified date for file %q. Error: %s", file.String(), lastModErr.Error())
			}

			file.Data(func(content io.ReadSeeker) error {
				http.ServeContent(w, r, filename, lastModifiedTime, content)
				return nil
			})

			return
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
