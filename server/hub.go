// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

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

func (hub *hub) ConnectionsByRoute(route string) []*connection {
	connectionsByRoute := make([]*connection, 0)

	for c := range h.connections {
		if c.Route == route {
			connectionsByRoute = append(connectionsByRoute, c)
		}
	}

	return connectionsByRoute
}

var h = hub{
	broadcast:   make(chan Message),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
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

				for _, c := range h.ConnectionsByRoute(m.Route) {
					select {
					case c.send <- m:
					default:
						delete(h.connections, c)
						close(c.send)
						go c.ws.Close()
					}
				}
			}
		}
	}
}
