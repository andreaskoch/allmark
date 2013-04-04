package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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
