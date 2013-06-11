// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	pather "github.com/andreaskoch/allmark/path"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var itemHandler = func(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path
	requestedPath = strings.TrimLeft(requestedPath, "/")

	fmt.Printf("Requested route %q\n", requestedPath)

	// make sure the request body is closed
	defer r.Body.Close()

	filePath, ok := routes[requestedPath]
	if !ok {

		// check for fallbacks before returning a 404
		if fallbackRoute, fallbackRouteFound := getFallbackRoute(requestedPath); fallbackRouteFound {
			redirect(w, r, fallbackRoute)
			return
		}

		error404Handler(w, r)
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		error404Handler(w, r)
		return
	}

	extension := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(extension)
	w.Header().Set("Content-Type", contentType)

	fmt.Fprintf(w, "%s", data)
}

func getFallbackRoute(requestedPath string) (fallbackRoute string, found bool) {

	requestedPath = pather.StripLeadingUrlDirectorySeperator(requestedPath)

	if len(requestedPath) == 0 {
		fmt.Printf("Fallback for %q is %q\n", requestedPath, pather.WebServerDefaultFilename)
		return pather.WebServerDefaultFilename, true
	}

	if strings.HasSuffix(requestedPath, pather.WebServerDefaultFilename) {
		fmt.Printf("No fallback found for %q\n", requestedPath)
		return "", false
	}

	route := pather.CombineUrlComponents(requestedPath, pather.WebServerDefaultFilename)
	if _, ok := routes[route]; ok {
		fmt.Printf("Fallback for %q is %q\n", requestedPath, route)
		return route, true
	}

	fmt.Printf("No fallback found for %q\n", requestedPath)
	return "", false
}
