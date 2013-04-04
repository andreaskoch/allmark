package server

import (
	"fmt"
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/renderer"
	"github.com/andreaskoch/docs/repository"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var routes map[string]repository.Pather

func Serve(repositoryPaths []string) {

	// An array of all indices for
	// the given repositories.
	indices := renderer.RenderRepositories(repositoryPaths)

	// Initialize the routing table
	initializeRoutes(indices)

	var error404Handler = func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path
		fmt.Fprintf(w, "Not found: %v", requestedPath)
	}

	var itemHandler = func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path

		fmt.Println(requestedPath)

		item, ok := routes[requestedPath]
		if !ok {

			// check for fallbacks before returning a 404
			if fallbackRoute, fallbackRouteFound := getFallbackRoute(requestedPath); fallbackRouteFound {
				redirect(w, r, fallbackRoute)
				return
			}

			error404Handler(w, r)
			return
		}

		data, err := ioutil.ReadFile(item.Path())
		if err != nil {
			error404Handler(w, r)
			return
		}

		fmt.Fprintf(w, "%s", data)
	}

	var indexDebugger = func(w http.ResponseWriter, r *http.Request) {
		for route, _ := range routes {
			fmt.Fprintln(w, route)
		}
	}

	// register handlers
	http.HandleFunc("/", itemHandler)
	http.HandleFunc("/debug/index", indexDebugger)

	// start http server
	http.ListenAndServe(":8080", nil)
}

func redirect(w http.ResponseWriter, r *http.Request, route string) {
	http.Redirect(w, r, route, http.StatusMovedPermanently)
}

func getFallbackRoute(requestedPath string) (fallbackRoute string, found bool) {

	if strings.HasSuffix(requestedPath, "/index.html") {
		return "", false
	}

	route := strings.TrimRight(requestedPath, "/") + "/" + "index.html"
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
