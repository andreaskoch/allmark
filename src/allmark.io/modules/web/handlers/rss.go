// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var itemsPerPage = 5

func RSS(headerWriter header.HeaderWriter,
	feedOrchestrator *orchestrator.FeedOrchestrator,
	templateProvider templates.Provider,
	error404Handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the current baseUrl
		baseUrl := getBaseUrlFromRequest(r)

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_XML)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the sitemap template
		feedTemplate, err := templateProvider.GetSubTemplate(baseUrl, templates.RssFeedTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// root entry / channel item
		rootEntry := feedOrchestrator.GetRootEntry(baseUrl)
		feedWrapper := renderFeedWrapper(feedTemplate, rootEntry)

		// get the sitemap content template
		feedContentTemplate, err := templateProvider.GetSubTemplate(baseUrl, templates.RssFeedContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		entries, found := feedOrchestrator.GetEntries(baseUrl, itemsPerPage, page)

		// display error 404 non-existing page has been requested
		if !found {
			error404Handler.ServeHTTP(w, r)
			return
		}

		sitemapContent := renderFeedEntries(feedContentTemplate, entries)

		feedWrapper = strings.Replace(feedWrapper, templates.ChildTemplatePlaceholder, sitemapContent, 1)

		fmt.Fprintf(w, "%s", feedWrapper)
	})
}

func renderFeedWrapper(templ *template.Template, root viewmodel.FeedEntry) string {
	buffer := new(bytes.Buffer)
	renderTemplate(templ, root, buffer)
	return buffer.String()
}

func renderFeedEntries(templ *template.Template, feedEntries []viewmodel.FeedEntry) string {

	rootCode := ""
	for _, entry := range feedEntries {
		buffer := new(bytes.Buffer)
		renderTemplate(templ, entry, buffer)
		rootCode += "\n" + buffer.String()
	}

	return rootCode
}
