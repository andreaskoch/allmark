// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"fmt"
	"net/http"
	"strings"
)

// Markdown returns a http handler which returns the markdown content of the requested item.
func Markdown(headerWriter header.HeaderWriter, viewModelOrchestrator *orchestrator.ViewModelOrchestrator, fallbackHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// strip the "md" or ".md" suffix from the path
		path := r.URL.Path
		path = strings.TrimSuffix(path, "markdown")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if viewModel, found := viewModelOrchestrator.GetFullViewModel(requestRoute); found {
			fmt.Fprintf(w, "%s", viewModel.Markdown)
			return
		}

		// fallback to the item handler
		fallbackHandler.ServeHTTP(w, r)
	})

}
