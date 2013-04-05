// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
)

var indexDebugger = func(w http.ResponseWriter, r *http.Request) {
	for route, _ := range routes {
		fmt.Fprintf(w, "%q => %q\n", route, routes[route])
	}
}
