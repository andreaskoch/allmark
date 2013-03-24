package server

import (
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/renderer"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var routes map[string]indexer.Addresser

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
			error404Handler(w, r)
			return
		}

		data, err := ioutil.ReadFile(item.GetAbsolutePath())
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

	http.HandleFunc("/", itemHandler)
	http.HandleFunc("/debug/index", indexDebugger)
	http.ListenAndServe(":8080", nil)
}

func initializeRoutes(indices []*indexer.Index) {

	routes = make(map[string]indexer.Addresser)

	for _, index := range indices {

		updateRouteTable := func(item *indexer.Item) {

			// get the item route and
			// add it to the routing table
			itemRoute := getHttpRouteFromFilePath(item.GetRelativePath(index.Path))
			registerRoute(itemRoute, item)

			// get the file routes and
			// add them to the routing table
			for _, file := range item.Files {
				fileRoute := getHttpRouteFromFilePath(file.GetRelativePath(index.Path))
				registerRoute(fileRoute, file)
			}
		}

		index.Walk(func(item *indexer.Item) {

			// add the current item to the route table
			updateRouteTable(item)

			// update route table again if item changes
			item.RegisterOnChangeCallback("UpdateRouteTableOnChange", func(i *indexer.Item) {
				i.IndexFiles()
				updateRouteTable(i)
			})
		})

	}
}

func getHttpRouteFromFilePath(path string) string {
	return strings.Replace(path, string(os.PathSeparator), "/", -1)
}

func registerRoute(route string, item indexer.Addresser) {

	if item == nil {
		log.Printf("Cannot add a route for an uninitialized item. Route: %#v\n", route)
		return
	}

	if strings.TrimSpace(route) == "" {
		log.Printf("Cannot add an empty route to the routing table. Item: %#v\n", item)
		return
	}

	routes[route] = item
}
