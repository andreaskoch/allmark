// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

func NewHub(logger logger.Logger, updateOrchestrator *orchestrator.UpdateOrchestrator) *Hub {
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

	updateOrchestrator *orchestrator.UpdateOrchestrator

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

	for connection := range hub.connections {

		if routeValue == connection.Route.Value() {
			connectionsByRoute = append(connectionsByRoute, connection)
		}
	}

	return connectionsByRoute
}

func (hub *Hub) run() {
	for {
		select {

		// subscribe a new connection
		case connection := <-hub.subscribe:
			{
				hub.logger.Debug("Subscring connection %s", connection.String())
				hub.logger.Debug("Number of Connections - Before: %v", len(hub.connections))

				// register the connection
				hub.connections[connection] = true

				hub.logger.Debug("Number of Connections - After: %v", len(hub.connections))
			}

		// unsubscribe an existing connection
		case connection := <-hub.unsubscribe:
			{
				hub.logger.Debug("Unsubscribing connection %s", connection.String())
				hub.logger.Debug("Number of Connections - Before: %v", len(hub.connections))

				// remove the connection
				delete(hub.connections, connection)

				hub.logger.Debug("Number of Connections - After: %v", len(hub.connections))
			}

		// handle broadcasts
		case broadcastMsg := <-hub.broadcast:
			{
				affectedConnections := hub.connectionsByRoute(broadcastMsg.Route)

				hub.logger.Debug("Revieved a broadcast message\n%#v", broadcastMsg)
				hub.logger.Debug("Connections affected: %v", len(affectedConnections))

				for _, connection := range affectedConnections {

					select {

					// send the message to the websocket
					case connection.send <- broadcastMsg:
						{
							hub.logger.Debug("Sending an update to: %s", connection.String())
						}

					default:
						{
							// todo: find out when this is happening
							hub.logger.Debug("Revieved a non-send message for %s", connection.String())
							delete(hub.connections, connection)
							go connection.ws.Close()
							hub.logger.Debug("Number of Connections: %v", len(hub.connections))
						}
					}

				}
			}
		}
	}
}
