// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"strings"
)

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan Message, 1),
		register:    make(chan *connection, 1),
		unregister:  make(chan *connection, 1),
		connections: make(map[*connection]bool),
	}
}

type Hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

func (hub *Hub) Message(viewModel viewmodel.Model) {
	hub.broadcast <- NewMessage(viewModel)
}

func (hub *Hub) Register(connection *connection) {
	hub.register <- connection
}

func (hub *Hub) UnRegister(connection *connection) {
	hub.unregister <- connection
}

func (hub *Hub) connectionsByRoute(route string) []*connection {
	connectionsByRoute := make([]*connection, 0)

	for c := range hub.connections {
		if strings.HasSuffix(route, c.Route) {
			connectionsByRoute = append(connectionsByRoute, c)
		}
	}

	return connectionsByRoute
}

func (hub *Hub) Run() {
	for {
		select {
		case c := <-hub.register:
			{
				hub.connections[c] = true
			}
		case c := <-hub.unregister:
			{
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
