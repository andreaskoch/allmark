// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"net/http"
	"strings"
)

func JSON(headerWriter header.HeaderWriter, viewModelOrchestrator *orchestrator.ViewModelOrchestrator, fallbackHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// strip the "json" or ".json" suffix from the path
		path := r.URL.Path
		path = strings.TrimSuffix(path, "json")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if viewModel, found := viewModelOrchestrator.GetFullViewModel(requestRoute); found {
			renderViewModelAsJSON(viewModel, w)
			return
		}

		// fallback to the item handler
		fallbackHandler.ServeHTTP(w, r)
	})

}
