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
	"os"
	"strings"
)

var (
	routes map[string]string

	useTempDir = true
)

const (

	// Dynamic Routes
	ItemHandlerRoute  = "/"
	DebugHandlerRoute = "/debug/index"

	// Static Routes
	ThemeFolderRoute = "/theme/"
)

func Serve(repositoryPath string) {

	index := renderer.RenderRepository(repositoryPath, useTempDir)

	// get the configuration
	config := config.GetConfig(repositoryPath)

	// Initialize the routing table
	initializeRoutes(index)

	// register handlers
	http.HandleFunc(ItemHandlerRoute, itemHandler)
	http.HandleFunc(DebugHandlerRoute, indexDebugger)

	// serve theme files
	if themeFolder := config.ThemeFolder(); util.DirectoryExists(themeFolder) {
		http.Handle(ThemeFolderRoute, http.StripPrefix(ThemeFolderRoute, http.FileServer(http.Dir(themeFolder))))
	}

	// start http server: http
	httpBinding := getHttpBinding(config)
	fmt.Printf("Starting http server %q\n", httpBinding)

	if err := http.ListenAndServe(httpBinding, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed with error: %v", err)
	}
}

func getHttpBinding(config *config.Config) string {

	// validate the port
	port := config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return fmt.Sprintf(":%v", port)
}

func initializeRoutes(index *repository.ItemIndex) {

	routes = make(map[string]string)

	pathProvider := path.NewProvider(index.Path(), useTempDir)
	for _, item := range index.Items() {
		registerItem(pathProvider, item)
	}
}

func registerItem(pathProvider *path.Provider, item *repository.Item) {

	// recurse for child items
	for _, child := range item.Childs() {
		registerItem(pathProvider, child)
	}

	// attach change listener
	item.OnChange("Update routing table on change", func(event *watcher.WatchEvent) {
		registerItem(pathProvider, item)
	})

	// get the item route and
	// add it to the routing table
	registerRoute(pathProvider, item)

	// get the file routes and
	// add them to the routing table
	for _, file := range item.Files.Items() {
		registerRoute(pathProvider, file)
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
