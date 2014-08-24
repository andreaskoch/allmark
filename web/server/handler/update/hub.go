// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

func NewHub(logger logger.Logger, updateOrchestrator orchestrator.UpdateOrchestrator) *Hub {
	hub := &Hub{
		logger: logger,

		updateOrchestrator: updateOrchestrator,

		broadcast:   make(chan Message, 1),
		subscribe:   make(chan *connection, 1),
		unsubscribe: make(chan *connection, 1),
		connections: make(map[*connection]bool),
	}

	// start the hub
	go hub.run()

	return hub
}

type Hub struct {
	logger logger.Logger

	updateOrchestrator orchestrator.UpdateOrchestrator

	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	subscribe chan *connection

	// Unsubscribe requests from connections.
	unsubscribe chan *connection
}

func (hub *Hub) Message(viewModel viewmodel.Model) {
	go func() {
		hub.logger.Debug("Broadcasting meesage %#v", viewModel)
		hub.broadcast <- NewMessage(viewModel)
	}()
}

func (hub *Hub) Subscribe(connection *connection) {
	hub.logger.Debug("Subscribing connection: %s", connection.String())

	// start watching for changes if there are no connections for this route
	if noConnectionsForRoute := len(hub.connectionsByRoute(connection.Route.Value())) == 0; noConnectionsForRoute {
		hub.updateOrchestrator.StartWatching(connection.Route)
	}

	hub.subscribe <- connection
}

func (hub *Hub) Unsubscribe(connection *connection) {

	// stop watching for changes if there are no more connections for this route
	if noConnectionsForRoute := len(hub.connectionsByRoute(connection.Route.Value())) <= 1; noConnectionsForRoute {
		hub.updateOrchestrator.StopWatching(connection.Route)
	}

	hub.unsubscribe <- connection
}

func (hub *Hub) connectionsByRoute(routeValue string) []*connection {
	connectionsByRoute := make([]*connection, 0)

	for c := range hub.connections {
		if routeValue == c.Route.Value() {
			connectionsByRoute = append(connectionsByRoute, c)
		}
	}

	return connectionsByRoute
}

func (hub *Hub) run() {
	for {
		select {
		case c := <-hub.subscribe:
			{
				hub.logger.Debug("Subscring connection %s", c.String())
				hub.logger.Debug("Number of Connections - Before: %v", len(hub.connections))
				hub.connections[c] = true
				hub.logger.Debug("Number of Connections - After: %v", len(hub.connections))
			}
		case c := <-hub.unsubscribe:
			{
				hub.logger.Debug("Unsubscribing connection %s", c.String())
				hub.logger.Debug("Number of Connections - Before: %v", len(hub.connections))
				delete(hub.connections, c)
				hub.logger.Debug("Number of Connections - After: %v", len(hub.connections))
			}
		case m := <-hub.broadcast:
			{
				hub.logger.Debug("Revieved a broadcast message %#v", m)

				affectedConnections := hub.connectionsByRoute(m.Route)
				for _, c := range affectedConnections {

					select {
					case c.send <- m:
						{
							hub.logger.Debug("Sending an update to: %s", c.String())
						}
					default:
						{
							// todo: find out when this is happening
							hub.logger.Debug("Revieved a non-send message for %s", c.String())
							delete(hub.connections, c)
							go c.ws.Close()
							hub.logger.Debug("Number of Connections: %v", len(hub.connections))
						}
					}

				}
			}
		}
	}
}
