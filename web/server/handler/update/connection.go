// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
)

func NewConnection(hub *Hub, ws *websocket.Conn, route route.Route) *connection {
	return &connection{
		Route: route,

		hub:  hub,
		send: make(chan Message, 10),
		ws:   ws,
	}
}

type connection struct {
	// The associated route.
	Route route.Route

	// the hub
	hub *Hub

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message
}

func (c *connection) String() string {
	return fmt.Sprintf("Connection (Route: %s, IP: %s)", c.Route.String(), c.ws.Request().RemoteAddr)
}

func (c *connection) Send(msg Message) {
	c.send <- msg
}

func (c *connection) Reader() {
	for {
		var message Message
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			break
		}

		c.hub.broadcast <- message
	}

	c.ws.Close()
	c.hub.Unsubscribe(c)
}

func (c *connection) Writer() {
	for message := range c.send {
		err := websocket.JSON.Send(c.ws, message)
		if err != nil {
			break
		}
	}

	c.ws.Close()
	c.hub.Unsubscribe(c)
}
