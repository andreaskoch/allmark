// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updatehandler

import (
	"code.google.com/p/go.net/websocket"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/update"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/gorilla/mux"
	"strings"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter, hub *update.Hub) *UpdateHandler {

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

	return &UpdateHandler{
		logger:                logger,
		config:                config,
		itemIndex:             itemIndex,
		hub:                   hub,
		patherFactory:         patherFactory,
		templateProvider:      templateProvider,
		error404Handler:       error404Handler,
		viewModelOrchestrator: viewModelOrchestrator,
	}
}

type UpdateHandler struct {
	logger logger.Logger
	config *config.Config

	itemIndex             *index.Index
	hub                   *update.Hub
	patherFactory         paths.PatherFactory
	templateProvider      *templates.Provider
	error404Handler       *errorhandler.ErrorHandler
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
}

func (handler *UpdateHandler) Func() func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {

		// get the path from the request variables
		vars := mux.Vars(ws.Request())
		path := vars["path"]

		// strip the "ws" or ".ws" suffix from the path
		path = strings.TrimSuffix(path, "ws")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// stage 1: check if there is a item for the request
		item, found := handler.itemIndex.IsMatch(*requestRoute)
		if !found {
			handler.logger.Debug("Route %q was not found.", requestRoute)
			return
		}

		// create a new connection
		c := update.NewConnection(handler.hub, ws, *requestRoute)

		// establish connection
		handler.logger.Debug("Establishing a connection for %q", requestRoute.String())
		handler.hub.Subscribe(c)

		// send an initial update
		go func() {
			// render the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := handler.viewModelOrchestrator.GetViewModel(pathProvider, item)
			c.Send(update.NewMessage(viewModel))
		}()

		defer func() {
			handler.hub.Unsubscribe(c)
		}()

		go c.Writer()

		c.Reader()
	}
}
