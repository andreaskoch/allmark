// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/renderer"
	"github.com/andreaskoch/allmark/repository"
	"log"
	"net/http"
	"strings"
)

var routes map[string]string

const (

	// Routes
	ItemHandlerRoute  = "/"
	DebugHandlerRoute = "/debug/index"
)

func Serve(repositoryPaths []string) {

	// An array of all indices for
	// the given repositories.
	indices := renderer.RenderRepositories(repositoryPaths)

	// Initialize the routing table
	initializeRoutes(indices)

	// register handlers
	http.HandleFunc(ItemHandlerRoute, itemHandler)
	http.HandleFunc(DebugHandlerRoute, indexDebugger)

	// start http server
	http.ListenAndServe(":8080", nil)
}

func getFallbackRoute(requestedPath string) (fallbackRoute string, found bool) {

	if strings.HasSuffix(requestedPath, path.WebServerDefaultFilename) {
		return "", false
	}

	route := path.CombineUrlComponents(requestedPath, path.WebServerDefaultFilename)
	if _, ok := routes[route]; ok {
		return route, true
	}

	return "", false
}

func initializeRoutes(indices []*repository.ItemIndex) {

	routes = make(map[string]string)

	for _, index := range indices {

		pathProvider := path.NewProvider(index.Path())

		updateRouteTable := func(item *repository.Item) {

			// get the item route and
			// add it to the routing table
			registerRoute(pathProvider, item)

			// get the file routes and
			// add them to the routing table
			for _, file := range item.Files.Items() {
				registerRoute(pathProvider, file)
			}
		}

		index.Walk(func(item *repository.Item) {

			// add the current item to the route table
			updateRouteTable(item)

			// update route table again if item changes
			item.RegisterOnChangeCallback("UpdateRouteTableOnChange", func(i *repository.Item) {
				updateRouteTable(i)
			})
		})

	}
}

func registerRoute(pathProvider *path.Provider, pather path.Pather) {

	if pather == nil {
		log.Printf("Cannot add a route for an uninitialized item %q.\n", pather.Path())
		return
	}

	route := pathProvider.GetWebRoute(pather)
	filePath := pathProvider.GetFilepath(pather)

	if strings.TrimSpace(route) == "" {
		log.Println("Cannot add an empty route to the routing table.")
		return
	}

	routes[route] = filePath
}
