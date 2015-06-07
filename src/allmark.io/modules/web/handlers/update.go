// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/handlers/update"
	"allmark.io/modules/web/header"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"strings"
)

func Update(logger logger.Logger,
	headerWriter header.HeaderWriter,
	updateOrchestrator *orchestrator.UpdateOrchestrator) websocket.Handler {

	hub := update.NewHub(logger, updateOrchestrator)

	updateChannel := make(chan orchestrator.Update, 1)
	updateOrchestrator.Subscribe(updateChannel)

	go func() {
		for update := range updateChannel {

			logger.Debug("Recieved an update for route %q", update.Route())

			// handle only modified items
			if update.Type() != orchestrator.UpdateTypeModified {
				continue
			}

			updatedModel, found := updateOrchestrator.GetUpdatedModel(update.Route())
			if !found {
				logger.Warn("The item for route %q was no longer found.", update.Route())
				return
			}

			hub.Message(updatedModel)
		}
	}()

	return websocket.Handler(func(ws *websocket.Conn) {

		// get the path from the request variables
		vars := mux.Vars(ws.Request())
		path := vars["path"]

		// strip the "ws" or ".ws" suffix from the path
		path = strings.TrimSuffix(path, "ws")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// stage 1: check if there is a item for the request
		if exists := updateOrchestrator.ItemExists(requestRoute); !exists {
			logger.Debug("Route %q was not found.", requestRoute)
			return
		}

		// create a new connection
		c := update.NewConnection(hub, ws, requestRoute)

		// establish connection
		logger.Debug("Establishing a connection for %q", requestRoute.String())
		hub.Subscribe(c)

		defer func() {
			hub.Unsubscribe(c)
		}()

		go c.Writer()

		c.Reader()
	})

}
