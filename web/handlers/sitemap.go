// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"fmt"
	"net/http"

	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/orchestrator"
	"github.com/andreaskoch/allmark/web/view/templates"
	"github.com/andreaskoch/allmark/web/view/viewmodel"

	"strings"
	"text/template"
)

func Sitemap(headerWriter header.HeaderWriter,
	navigationOrchestrator *orchestrator.NavigationOrchestrator,
	sitemapOrchestrator *orchestrator.SitemapOrchestrator,
	templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseURLFromRequest(r)

		// get the sitemap template
		sitemapTemplate, err := templateProvider.GetSitemapTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// get the sitemap-entry template
		sitemapEntryTemplate, childPlaceholder, err := templateProvider.GetSitemapEntryTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page parameters
		pageTitle := "Sitemap"
		pageType := "sitemap"
		descriptionText := "A list of all items in this repository."

		// Page model
		viewModel := viewmodel.Model{}
		viewModel.Type = pageType
		viewModel.Title = pageTitle
		viewModel.PageTitle = sitemapOrchestrator.GetPageTitle(pageTitle)
		viewModel.Description = descriptionText
		viewModel.ToplevelNavigation = navigationOrchestrator.GetToplevelNavigation()
		viewModel.BreadcrumbNavigation = navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		sitemapPageModel := viewmodel.Sitemap{}
		sitemapPageModel.Model = viewModel
		sitemapPageModel.Tree = renderSitemapEntryTemplate(sitemapEntryTemplate, sitemapOrchestrator.GetSitemap(), childPlaceholder)

		renderTemplate(sitemapTemplate, sitemapPageModel, w)
	})
}

func renderSitemapEntryTemplate(template *template.Template, entry viewmodel.SitemapEntry, childPlaceholder string) string {
	content, err := getRenderedCode(template, entry)
	if err != nil {
		return err.Error()
	}

	childCode := ""
	for _, child := range entry.Children {
		childCode += renderSitemapEntryTemplate(template, child, childPlaceholder)
	}

	content = strings.Replace(content, childPlaceholder, childCode, -1)

	return content
}
