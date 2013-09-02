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
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

var (
	items  = repository.Items{}
	routes = make(map[string]*Route)

	useTempDir = true

	error404Handler   func(w http.ResponseWriter, r *http.Request)
	xmlSitemapHandler func(w http.ResponseWriter, r *http.Request)
	tagmapHandler     func(w http.ResponseWriter, r *http.Request)
	sitemapHandler    func(w http.ResponseWriter, r *http.Request)
	rssHandler        func(w http.ResponseWriter, r *http.Request)
)

const (

	// Dynamic Routes
	ItemHandlerRoute       = "/"
	TagmapHandlerRoute     = "/tags.html"
	SitemapHandlerRoute    = "/sitemap.html"
	XmlSitemapHandlerRoute = "/sitemap.xml"
	RssHandlerRoute        = "/rss.xml"
	RobotsTxtHandlerRoute  = "/robots.txt"
	DebugHandlerRoute      = "/debug/index"
	WebSocketHandlerRoute  = "/ws"

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

	// start the renderer
	server.renderer.Execute()

	// initialize the 404 handler
	error404Handler = func(w http.ResponseWriter, r *http.Request) {

		// set 404 status code
		w.WriteHeader(http.StatusNotFound)

		// write 404 page
		server.renderer.Error404(w)
	}

	// initialize the xml sitemap handler
	xmlSitemapHandler = func(w http.ResponseWriter, r *http.Request) {
		server.renderer.XMLSitemap(w, getHostnameFromRequest(r))
	}

	// initialize the tagmap handler
	tagmapHandler = func(w http.ResponseWriter, r *http.Request) {
		server.renderer.Tags(w, getHostnameFromRequest(r))
	}

	// initialize the sitemap handler
	sitemapHandler = func(w http.ResponseWriter, r *http.Request) {
		server.renderer.Sitemap(w, getHostnameFromRequest(r))
	}

	// initialize the RSS handler
	rssHandler = func(w http.ResponseWriter, r *http.Request) {
		server.renderer.RSS(w, getHostnameFromRequest(r))
	}

	// start a change listener
	server.listenForChanges()

	// start the websocket hub
	go h.run()

	// register handlers
	http.HandleFunc(ItemHandlerRoute, itemHandler)
	http.HandleFunc(TagmapHandlerRoute, tagmapHandler)
	http.HandleFunc(SitemapHandlerRoute, sitemapHandler)
	http.HandleFunc(XmlSitemapHandlerRoute, xmlSitemapHandler)
	http.HandleFunc(RssHandlerRoute, rssHandler)
	http.HandleFunc(RobotsTxtHandlerRoute, robotsTxtHandler)
	http.HandleFunc(DebugHandlerRoute, indexDebugger)
	http.Handle(WebSocketHandlerRoute, websocket.Handler(webSocketHandler))

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

func (server *Server) listenForChanges() {

	go func() {

		for {
			select {
			case item := <-server.renderer.Rendered:
				if item != nil {
					// register the item
					server.registerItem(item)

					// send update event to connected browsers
					h.broadcast <- UpdateMessage(item.Model)
				}
			case item := <-server.renderer.Removed:
				if item != nil {
					// un-register the item
					server.unregisterItem(item)
				}
			}
		}
	}()

}

func (server *Server) unregisterItem(item *repository.Item) {

	// recurse for child items
	for _, child := range item.Childs {
		server.unregisterItem(child)
	}

	// add item to list
	newItemList := repository.Items{}
	for _, entry := range items {
		if entry.String() != item.String() {
			newItemList = append(newItemList, entry)
		}
	}

	items = newItemList

	// unregister the item
	server.unregisterRoute(item)

	// unregister all item files
	for _, file := range item.Files.Items() {
		server.unregisterRoute(file)
	}
}

func (server *Server) registerItem(item *repository.Item) {

	// recurse for child items
	for _, child := range item.Childs {
		server.registerItem(child)
	}

	// add item to list
	items = append(items, item)

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
		log.Println("Cannot add a route for an uninitialized item.")
		return
	}

	route, err := server.getRoute(pather)
	if err != nil {
		log.Println("Route could not be registered. Error: %s", err)
		return
	}

	routes[route.String()] = route
}

func (server *Server) unregisterRoute(pather path.Pather) {

	if pather == nil {
		log.Println("Cannot unregister a route for an uninitialized item.")
		return
	}

	route, err := server.getRoute(pather)
	if err != nil {
		log.Println("Route could not be un-registered. Error: %s", err)
		return
	}

	delete(routes, route.String())
}

func (server *Server) getRoute(pather path.Pather) (*Route, error) {
	route := server.pathProvider.GetWebRoute(pather)
	filepath := server.pathProvider.GetFilepath(pather)
	return newRoute(route, filepath)
}

func getRequestedPathFromRequest(r *http.Request) string {
	requestedPath := r.URL.Path
	requestedPath = strings.TrimLeft(requestedPath, "/")
	requestedPath = util.EncodeUrl(requestedPath)
	return requestedPath
}

func getHostnameFromRequest(r *http.Request) string {
	return r.Host
}

func normalizeRoute(route string) string {
	return strings.ToLower(util.DecodeUrl(route))
}
