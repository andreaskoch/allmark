// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"net/url"
)

func webSocketHandler(ws *websocket.Conn) {

	// get the request uri
	request := ws.Request()
	requestUri, err := url.Parse(request.RequestURI)
	if err != nil {

		//  close the connection
		fmt.Printf("Cannot establish a web socket connection without the request url.\nError:\n%s\n", err)
		ws.Close()
		return

	}

	// extract the route parameter from the request uri
	queryParameters := requestUri.Query()
	routeParam := queryParameters.Get("route")
	if routeParam == "" {

		//  close the connection
		ws.Close()
		return

	}

	// check if the route exists
	routeParam = path.StripLeadingUrlDirectorySeperator(routeParam)
	routeParam = normalizeRoute(routeParam)

	route, routeExists := routes[routeParam]
	if !routeExists {
		fmt.Printf("Cannot establish a web socket connection. The route %q does not exist.\n", routeParam)
		return
	}

	// create a new connection
	c := &connection{
		Route: route.Original(),
		send:  make(chan Message, 10),
		ws:    ws,
	}

	// establish connection
	fmt.Printf("Establishing a connection for %q\n", route.Original())
	h.register <- c

	defer func() {
		h.unregister <- c
	}()

	go c.writer()

	c.reader()
}
