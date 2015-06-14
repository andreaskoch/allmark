// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"github.com/gorilla/mux"
	"net/http"
)

// AliasLookup creates a http handler which redirects aliases to their documents.
func AliasLookup(headerWriter header.HeaderWriter, viewModelOrchestrator *orchestrator.ViewModelOrchestrator, fallbackHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		alias := vars["alias"]

		// locate the correct viewmodel for the given alias
		viewModel, found := viewModelOrchestrator.GetViewModelByAlias(alias)
		if !found {
			// no model found for the alias -> use fallback handler
			fallbackHandler.ServeHTTP(w, r)
			return
		}

		// determine the redirect url
		baseUrl := getBaseUrlFromRequest(r)
		redirectUrl := baseUrl + viewModel.BaseUrl
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
	})

}
