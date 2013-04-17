// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"io/ioutil"
	"net/http"
	"strings"
)

var itemHandler = func(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path
	requestedPath = strings.TrimLeft(requestedPath, path.UrlDirectorySeperator)

	fmt.Printf("Requested route %q\n", requestedPath)

	filePath, ok := routes[requestedPath]
	if !ok {

		error404Handler(w, r)
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		error404Handler(w, r)
		return
	}

	fmt.Fprintf(w, "%s", data)
}
