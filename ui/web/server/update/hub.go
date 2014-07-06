// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewHub(logger logger.Logger, updateHub dataaccess.UpdateHub) *Hub {
	return &Hub{
		logger: logger,

		updateHub: updateHub,

		broadcast:   make(chan Message, 1),
		subscribe:   make(chan *connection, 1),
		unsubscribe: make(chan *connection, 1),
		connections: make(map[*connection]bool),
	}
}

type Hub struct {
	logger logger.Logger

	updateHub dataaccess.UpdateHub

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
		hub.broadcast <- NewMessage(viewModel)
	}()
}

func (hub *Hub) Subscribe(connection *connection) {
	go hub.updateHub.StartWatching(connection.Route)
	hub.subscribe <- connection
}

func (hub *Hub) Unsubscribe(connection *connection) {
	go hub.updateHub.StopWatching(connection.Route)

	hub.unsubscribe <- connection
}

func (hub *Hub) connectionsByRoute(routeValue string) []*connection {
	connectionsByRoute := make([]*connection, 0)

	hub.logger.Debug("Number of Connections %s", len(hub.connections))

	for c := range hub.connections {
		if routeValue == c.Route.Value() {
			connectionsByRoute = append(connectionsByRoute, c)
		}
	}

	return connectionsByRoute
}

func (hub *Hub) Run() {
	for {
		select {
		case c := <-hub.subscribe:
			{
				hub.logger.Debug("üüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüüü")
				hub.logger.Debug("Subscriping connection %s", c.Route.Value())
				hub.connections[c] = true
			}
		case c := <-hub.unsubscribe:
			{
				hub.logger.Debug("Unsubscriping connection %s", c.Route.Value())
				delete(hub.connections, c)
				close(c.send)
			}
		case m := <-hub.broadcast:
			{
				affectedConnections := hub.connectionsByRoute(m.Route)
				for _, c := range affectedConnections {

					select {
					case c.send <- m:
					default:
						delete(hub.connections, c)

						// todo: introduce a maanger which sends a signal if a route is removed and closes the channel
						// if I just call close there this will fail quite often if the channel has already been closed.
						//close(c.send)

						go c.ws.Close()
					}

				}
			}
		}
	}
}
