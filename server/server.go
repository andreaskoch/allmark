package server

import (
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/renderer"
	"github.com/andreaskoch/docs/repository"
	"log"
	"net/http"
	"strings"
)

var routes map[string]repository.Pather

const (
	UrlDirectorySeperator = "/"
	DefaultFile           = "index.html"

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

	if strings.HasSuffix(requestedPath, "/index.html") {
		return "", false
	}

	route := CombineUrlComponents(requestedPath, "index.html")
	if _, ok := routes[route]; ok {
		return route, true
	}

	return "", false
}

func initializeRoutes(indices []*repository.ItemIndex) {

	routes = make(map[string]repository.Pather)

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

func registerRoute(pathProvider *path.Provider, pather repository.Pather) {

	if pather == nil {
		log.Printf("Cannot add a route for an uninitialized item %q.\n", pather.Path())
		return
	}

	route := pathProvider.GetWebRoute(pather)

	if strings.TrimSpace(route) == "" {
		log.Printf("Cannot add an empty route to the routing table. Item %q\n", pather.Path())
		return
	}

	routes[route] = pather
}
