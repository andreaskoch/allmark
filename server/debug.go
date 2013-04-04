package server

import (
	"fmt"
	"net/http"
)

var indexDebugger = func(w http.ResponseWriter, r *http.Request) {
	for route, _ := range routes {
		fmt.Fprintln(w, route)
	}
}
