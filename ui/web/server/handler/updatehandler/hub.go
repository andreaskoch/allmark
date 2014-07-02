// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updatehandler

import (
	"strings"
)

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan Message, 1),
	register:    make(chan *connection, 1),
	unregister:  make(chan *connection, 1),
	connections: make(map[*connection]bool),
}

func (hub *hub) ConnectionsByRoute(route string) []*connection {
	connectionsByRoute := make([]*connection, 0)

	for c := range h.connections {
		if strings.HasSuffix(route, c.Route) {
			connectionsByRoute = append(connectionsByRoute, c)
		}
	}

	return connectionsByRoute
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			{
				h.connections[c] = true
			}
		case c := <-h.unregister:
			{
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			{
				affectedConnections := h.ConnectionsByRoute(m.Route)
				for _, c := range affectedConnections {

					select {
					case c.send <- m:
					default:
						delete(h.connections, c)

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
