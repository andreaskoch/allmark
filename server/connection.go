// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message
}

func (c *connection) reader() {
	for {
		var message Message
		err := websocket.JSON.Receive(c.ws, &message)
		if err != nil {
			break
		}
		h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := websocket.JSON.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
	fmt.Printf("%#v\n", ws.Request())

	c := &connection{send: make(chan Message, 256), ws: ws}

	h.register <- c

	defer func() {
		h.unregister <- c
	}()

	go c.writer()

	c.reader()
}
