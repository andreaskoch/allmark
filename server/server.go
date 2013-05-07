// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"code.google.com/p/go.net/websocket"
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
	ItemHandlerRoute      = "/"
	DebugHandlerRoute     = "/debug/index"
	WebSocketHandlerRoute = "/ws"

	// Static Routes
	ThemeFolderRoute = "/theme/"
)

type Server struct {
	repositoryPath string
	pathProvider   *path.Provider
	config         *config.Config
	renderer       *renderer.Renderer
}

func New(repositoryPath string, config *config.Config, useTempDir bool) *Server {

	return &Server{
		repositoryPath: repositoryPath,
		pathProvider:   path.NewProvider(repositoryPath, useTempDir),
		config:         config,
		renderer:       renderer.New(repositoryPath, config, useTempDir),
	}

}

func (server *Server) Serve() {

	index := server.renderer.Execute()

	// Initialize the routing table
	server.initializeRoutes(index)

	// start the websocket hub
	go h.run()

	// register handlers
	http.HandleFunc(ItemHandlerRoute, itemHandler)
	http.HandleFunc(DebugHandlerRoute, indexDebugger)
	http.Handle(WebSocketHandlerRoute, websocket.Handler(wsHandler))

	// serve theme files
	if themeFolder := server.config.ThemeFolder(); util.DirectoryExists(themeFolder) {
		http.Handle(ThemeFolderRoute, http.StripPrefix(ThemeFolderRoute, http.FileServer(http.Dir(themeFolder))))
	}

	// start http server: http
	httpBinding := server.getHttpBinding()
	fmt.Printf("Starting http server %q\n", httpBinding)

	if err := http.ListenAndServe(httpBinding, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed with error: %v", err)
	}
}

func (server *Server) getHttpBinding() string {

	// validate the port
	port := server.config.Server.Http.Port
	if port < 1 || port > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", port, 1, math.MaxUint16))
	}

	return fmt.Sprintf(":%v", port)
}

func (server *Server) initializeRoutes(index *repository.ItemIndex) {

	routes = make(map[string]string)

	for _, item := range index.Items() {
		server.registerItem(item)
	}
}

func (server *Server) registerItem(item *repository.Item) {

	// recurse for child items
	for _, child := range item.Childs() {
		server.registerItem(child)
	}

	// attach change listener
	item.OnChange("Update routing table on change", func(event *watcher.WatchEvent) {

		// re-register item on change
		server.registerItem(item)

		// send update event to connected browsers
		h.broadcast <- UpdateMessage(item.ViewModel)

	})

	// get the item route and
	// add it to the routing table
	server.registerRoute(item)

	// get the file routes and
	// add them to the routing table
	for _, file := range item.Files.Items() {
		server.registerRoute(file)
	}
}

func (server *Server) registerRoute(pather path.Pather) {

	if pather == nil {
		log.Printf("Cannot add a route for an uninitialized item %q.\n", pather.Path())
		return
	}

	route := server.pathProvider.GetWebRoute(pather)
	filePath := server.pathProvider.GetFilepath(pather)

	if strings.TrimSpace(route) == "" {
		log.Println("Cannot add an empty route to the routing table.")
		return
	}

	routes[route] = filePath
}
