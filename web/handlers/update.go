// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/web/handlers/update"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"github.com/elWyatt/allmark/web/view/templates"
	"github.com/elWyatt/allmark/web/view/templates/templatenames"
	"github.com/elWyatt/allmark/web/view/viewmodel"
	"golang.org/x/net/websocket"
	"strings"
)

func Update(logger logger.Logger,
	headerWriter header.HeaderWriter,
	templateProvider templates.Provider,
	updateOrchestrator *orchestrator.UpdateOrchestrator) websocket.Handler {

	hub := update.NewHub(logger, updateOrchestrator)

	updateChannel := make(chan orchestrator.Update, 1)
	updateOrchestrator.Subscribe(updateChannel)

	go func() {
		for update := range updateChannel {

			logger.Info("Received an update for route %q: %s", update.Route(), update.String())

			// handle new or modified items
			if update.Type() == orchestrator.UpdateTypeNew || update.Type() == orchestrator.UpdateTypeModified {

				// send the latest viewmodel to the client
				viewModel, found := updateOrchestrator.GetUpdatedModel(update.Route())
				if !found {
					logger.Warn("The item for route %q was no longer found.", update.Route())
					return
				}

				var updateModel viewmodel.Update
				updateModel.Model = viewModel

				snippets := make(map[string]string)
				snippets["aliases"] = renderSnippet(templateProvider, templatenames.Aliases, viewModel)
				snippets["tags"] = renderSnippet(templateProvider, templatenames.Tags, viewModel)
				snippets["publisher"] = renderSnippet(templateProvider, templatenames.Publisher, viewModel)
				snippets["toplevelnavigation"] = renderSnippet(templateProvider, templatenames.ToplevelNavigation, viewModel)
				snippets["breadcrumbnavigation"] = renderSnippet(templateProvider, templatenames.BreadcrumbNavigation, viewModel)
				snippets["itemnavigation"] = renderSnippet(templateProvider, templatenames.ItemNavigation, viewModel)
				snippets["children"] = renderSnippet(templateProvider, templatenames.Children, viewModel)
				snippets["tagcloud"] = renderSnippet(templateProvider, templatenames.TagCloud, viewModel)

				updateModel.Snippets = snippets

				hub.Message(updateModel)

			} else if update.Type() == orchestrator.UpdateTypeNew {

				// send an empty update to the client
				hub.Message(viewmodel.Update{})

			} else {

				logger.Debug("No action for update %s", update.String())

			}

		}
	}()

	return websocket.Handler(func(ws *websocket.Conn) {

		// strip the "ws" or ".ws" suffix from the path
		path := ws.Request().URL.Path
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

func renderSnippet(templateProvider templates.Provider, templateName string, viewmodel interface{}) string {

	// get the search result content template
	hostname := "" // TODO: get real hostname
	subTemplate, err := templateProvider.GetSnippetTemplate(templateName, hostname)
	if err != nil {
		return err.Error()
	}

	code, err := getRenderedCode(subTemplate, viewmodel)
	if err != nil {
		return err.Error()
	}

	return code
}
