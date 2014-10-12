// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"net/http"
	"strings"
	"text/template"
)

var itemsPerPage = 5

type Rss struct {
	logger logger.Logger

	templateProvider templates.Provider
	error404Handler  Handler
	feedOrchestrator orchestrator.FeedOrchestrator
}

func (handler *Rss) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/xml")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the sitemap template
		feedTemplate, err := handler.templateProvider.GetSubTemplate(templates.RssFeedTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// get the current hostname
		hostname := getHostnameFromRequest(r)

		// root entry / channel item
		rootEntry := handler.feedOrchestrator.GetRootEntry(hostname)
		feedWrapper := renderFeedWrapper(feedTemplate, rootEntry)

		// get the sitemap content template
		feedContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.RssFeedContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		entries := handler.feedOrchestrator.GetEntries(hostname, itemsPerPage, page)

		// display error 404 non-existing page has been requested
		if page > 1 && len(entries) == 0 {
			error404Handler := handler.error404Handler.Func()
			error404Handler(w, r)
			return
		}

		sitemapContent := renderFeedEntries(feedContentTemplate, entries)

		feedWrapper = strings.Replace(feedWrapper, templates.ChildTemplatePlaceholder, sitemapContent, 1)

		fmt.Fprintf(w, "%s", feedWrapper)
	}
}

func renderFeedWrapper(templ *template.Template, root viewmodel.FeedEntry) string {
	buffer := new(bytes.Buffer)
	renderTemplate(root, templ, buffer)
	return buffer.String()
}

func renderFeedEntries(templ *template.Template, feedEntries []viewmodel.FeedEntry) string {

	rootCode := ""
	for _, entry := range feedEntries {
		buffer := new(bytes.Buffer)
		renderTemplate(entry, templ, buffer)
		rootCode += "\n" + buffer.String()
	}

	return rootCode
}
