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
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/gorilla/mux"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *UpdateHandler {

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

	// start the websocket hub
	go h.run()

	return &UpdateHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		config:                config,
		patherFactory:         patherFactory,
		templateProvider:      templateProvider,
		error404Handler:       error404Handler,
		viewModelOrchestrator: viewModelOrchestrator,
	}
}

type UpdateHandler struct {
	logger                logger.Logger
	itemIndex             *index.Index
	config                *config.Config
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

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// stage 1: check if there is a item for the request
		_, found := handler.itemIndex.IsMatch(*requestRoute)
		if !found {
			handler.logger.Debug("Route %q was not found.", requestRoute)
			return
		}

		// create a new connection
		c := &connection{
			Route: requestRoute.Value(),
			send:  make(chan Message, 10),
			ws:    ws,
		}

		// establish connection
		handler.logger.Debug("Establishing a connection for %q", requestRoute.String())
		h.register <- c

		defer func() {
			h.unregister <- c
		}()

		go c.writer()

		c.reader()
	}
}
