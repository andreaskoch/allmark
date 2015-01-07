// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/handler/update"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"strings"
)

type Update struct {
	logger logger.Logger

	updateOrchestrator *orchestrator.UpdateOrchestrator
}

func (handler *Update) Func() func(ws *websocket.Conn) {

	hub := update.NewHub(handler.logger, handler.updateOrchestrator)

	// attach a update-handler
	handler.updateOrchestrator.OnUpdate(func(updatedRoute route.Route) {
		updatedModel, found := handler.updateOrchestrator.GetUpdatedModel(updatedRoute)
		if !found {
			handler.logger.Warn("The item for route %q was no longer found.", updatedRoute)
			return
		}

		hub.Message(updatedModel)
	})

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
		if exists := handler.updateOrchestrator.ItemExists(requestRoute); !exists {
			handler.logger.Debug("Route %q was not found.", requestRoute)
			return
		}

		// create a new connection
		c := update.NewConnection(hub, ws, requestRoute)

		// establish connection
		handler.logger.Debug("Establishing a connection for %q", requestRoute.String())
		hub.Subscribe(c)

		defer func() {
			hub.Unsubscribe(c)
		}()

		go c.Writer()

		c.Reader()
	}
}
