// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var itemsPerPage = 5

type Rss struct {
	logger           logger.Logger
	headerWriter     header.HeaderWriter
	templateProvider templates.Provider
	error404Handler  Handler
	feedOrchestrator *orchestrator.FeedOrchestrator
}

func (handler *Rss) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the current hostname
		hostname := getBaseUrlFromRequest(r)

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_XML)

		// read the page url-parameter
		page, pageParameterIsAvailable := getPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the sitemap template
		feedTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.RssFeedTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// root entry / channel item
		rootEntry := handler.feedOrchestrator.GetRootEntry(hostname)
		feedWrapper := renderFeedWrapper(feedTemplate, rootEntry)

		// get the sitemap content template
		feedContentTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.RssFeedContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		entries, found := handler.feedOrchestrator.GetEntries(hostname, itemsPerPage, page)

		// display error 404 non-existing page has been requested
		if !found {
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
