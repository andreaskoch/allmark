// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"net/http"
	"os"
	"strings"
)

var itemHandler = func(w http.ResponseWriter, r *http.Request) {
	requestedPath := getRequestedPathFromRequest(r)
	fmt.Printf("Requested route %q\n", requestedPath)

	// make sure the request body is closed
	defer r.Body.Close()

	// check if the route is registered
	normalizedRequestPath := normalizeRoute(requestedPath)
	route, ok := routes[normalizedRequestPath]
	if !ok {

		// check for fallbacks before returning a 404
		if fallbackRoute, fallbackRouteFound := getFallbackRoute(normalizedRequestPath); fallbackRouteFound {

			// redirect to the fallback route
			redirect(w, r, fallbackRoute)
			return
		}

		error404Handler(w, r)
		return
	}

	// check the casing
	if requestedPath != route.Original() {

		fmt.Println("Requested path", requestedPath)
		fmt.Println("Original path", route.Original())

		// redirect to the route with the correct casing
		redirect(w, r, route.Original())
		return
	}

	// open the file
	file, err := os.Open(route.Filepath())
	if err != nil {
		error404Handler(w, r)
		return
	}

	defer file.Close()

	// get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		error404Handler(w, r)
		return
	}

	// serve the file
	http.ServeContent(w, r, route.Filepath(), fileInfo.ModTime(), file)
}

func getFallbackRoute(requestedPath string) (fallbackRoute string, found bool) {

	// empty path
	if len(path.StripLeadingUrlDirectorySeperator(requestedPath)) == 0 {
		fmt.Printf("Fallback for %q is %q\n", requestedPath, path.WebServerDefaultFilename)
		return path.WebServerDefaultFilename, true
	}

	// index.html has already been requested
	if strings.HasSuffix(requestedPath, path.WebServerDefaultFilename) {
		fmt.Printf("No fallback found for %q\n", requestedPath)
		return "", false
	}

	// try to add index.html
	route := path.CombineUrlComponents(requestedPath, path.WebServerDefaultFilename)
	if _, ok := routes[route]; ok {
		fmt.Printf("Fallback for %q is %q\n", requestedPath, route)
		return route, true
	}

	// no route found
	fmt.Printf("No fallback found for %q\n", requestedPath)
	return "", false
}
