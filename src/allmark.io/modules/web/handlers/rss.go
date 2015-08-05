// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"fmt"
	"net/http"
)

var itemsPerPage = 5

func RSS(headerWriter header.HeaderWriter,
	feedOrchestrator *orchestrator.FeedOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the current baseURL
		hostname := getBaseURLFromRequest(r)

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_XML)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromURL(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the RSS template
		feedTemplate, err := templateProvider.GetRSSTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		feedModel, err := feedOrchestrator.GetFeed(hostname, itemsPerPage, page)

		// display error 404 non-existing page has been requested
		if err != nil {
			error404Handler.ServeHTTP(w, r)
			return
		}

		renderTemplate(feedTemplate, feedModel, w)
	})
}
