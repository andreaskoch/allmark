// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/renderer"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"github.com/andreaskoch/allmark/watcher"
	"log"
	"math"
	"net/http"
	"strings"
)

var routes map[string]string

const (

	// Routes
	ItemHandlerRoute  = "/"
	DebugHandlerRoute = "/debug/index"
)

func Serve(repositoryPath string) {

	index := renderer.RenderRepository(repositoryPath)

	// get the configuration
	config := config.GetConfig(repositoryPath)

	// Initialize the routing table
	initializeRoutes(index)

	// register handlers
	http.HandleFunc(ItemHandlerRoute, itemHandler)
	http.HandleFunc(DebugHandlerRoute, indexDebugger)

	// serve theme files
	if themeFolder := config.ThemeFolder(); util.DirectoryExists(themeFolder) {
		staticFolder := "/theme/"
		http.Handle(staticFolder, http.StripPrefix(staticFolder, http.FileServer(http.Dir(themeFolder))))
	}

	// start http server: http
	httpBinding := getHttpBinding(config)
	fmt.Printf("Starting http server %q\n", httpBinding)

	http.ListenAndServe(httpBinding, nil)
}

func getHttpBinding(config *config.Config) string {

	// validate the port
	port := config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return fmt.Sprintf(":%v", port)
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

func initializeRoutes(index *repository.ItemIndex) {

	routes = make(map[string]string)

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
		item.OnChange("Update routing table on change", func(event *watcher.WatchEvent) {
			updateRouteTable(item)
		})
	})
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
