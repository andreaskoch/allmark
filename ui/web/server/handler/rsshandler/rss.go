// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsshandler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
	"strings"
	"text/template"
)

var itemsPerPage = 5

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) *RssHandler {

	templateProvider := templates.NewProvider(".")
	feedOrchestrator := orchestrator.NewFeedOrchestrator(itemIndex, converter)

	return &RssHandler{
		logger:           logger,
		itemIndex:        itemIndex,
		config:           config,
		patherFactory:    patherFactory,
		templateProvider: templateProvider,
		feedOrchestrator: feedOrchestrator,
	}
}

type RssHandler struct {
	logger           logger.Logger
	itemIndex        *index.ItemIndex
	config           *config.Config
	patherFactory    paths.PatherFactory
	templateProvider *templates.Provider
	feedOrchestrator orchestrator.FeedOrchestrator
}

func (handler *RssHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// read the page url-parameter
		page, pageParameterIsAvailable := handlerutil.GetPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		// get the sitemap template
		xmlSitemapTemplate, err := handler.templateProvider.GetSubTemplate(templates.RssFeedTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		// root entry / channel item
		rootEntry := handler.feedOrchestrator.GetRootEntry(pathProvider)
		feedWrapper := renderFeedWrapper(xmlSitemapTemplate, rootEntry)

		// get the sitemap content template
		feedContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.RssFeedContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		entries := handler.feedOrchestrator.GetEntries(pathProvider, itemsPerPage, page)

		sitemapContent := renderFeedEntries(feedContentTemplate, entries)

		feedWrapper = strings.Replace(feedWrapper, templates.ChildTemplatePlaceholder, sitemapContent, 1)

		fmt.Fprintf(w, "%s", feedWrapper)
	}
}

func renderFeedWrapper(templ *template.Template, root viewmodel.FeedEntry) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(root, templ, buffer)
	return buffer.String()
}

func renderFeedEntries(templ *template.Template, feedEntries []viewmodel.FeedEntry) string {

	rootCode := ""
	for _, entry := range feedEntries {
		buffer := new(bytes.Buffer)
		handlerutil.RenderTemplate(entry, templ, buffer)
		rootCode += "\n" + buffer.String()
	}

	return rootCode
}
