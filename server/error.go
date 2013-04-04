package server

import (
	"fmt"
	"net/http"
)

var error404Handler = func(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Path
	fmt.Fprintf(w, "Not found: %v", requestedPath)
}
