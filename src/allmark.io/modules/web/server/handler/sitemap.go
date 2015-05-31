// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
)

type Sitemap struct {
	logger                 logger.Logger
	headerWriter           header.HeaderWriter
	navigationOrchestrator *orchestrator.NavigationOrchestrator
	sitemapOrchestrator    *orchestrator.SitemapOrchestrator
	templateProvider       templates.Provider
}

func (handler *Sitemap) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_HTML)

		hostname := getBaseUrlFromRequest(r)

		// get the sitemap content template
		sitemapContentTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.SitemapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// get the sitemap template
		sitemapTemplate, err := handler.templateProvider.GetFullTemplate(hostname, templates.SitemapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// Page parameters
		pageTitle := "Sitemap"
		pageType := "sitemap"
		descriptionText := "A list of all items in this repository."

		// Page content
		sitemapContentModel := handler.sitemapOrchestrator.GetSitemap()
		sitemapContent := renderSitemapEntry(sitemapContentTemplate, sitemapContentModel)

		// Page model
		sitemapPageModel := viewmodel.Model{
			Content: sitemapContent,
		}

		sitemapPageModel.Type = pageType
		sitemapPageModel.Title = pageTitle
		sitemapPageModel.PageTitle = handler.sitemapOrchestrator.GetPageTitle(pageTitle)
		sitemapPageModel.Description = descriptionText
		sitemapPageModel.ToplevelNavigation = handler.navigationOrchestrator.GetToplevelNavigation()
		sitemapPageModel.BreadcrumbNavigation = handler.navigationOrchestrator.GetBreadcrumbNavigation(route.New())

		renderTemplate(sitemapPageModel, sitemapTemplate, w)
	}
}

func renderSitemapEntry(templ *template.Template, sitemapModel viewmodel.Sitemap) string {

	// render
	buffer := new(bytes.Buffer)
	renderTemplate(sitemapModel, templ, buffer)

	// get the produced html code
	rootCode := buffer.String()

	if len(sitemapModel.Childs) > 0 {

		// render all childs
		childCode := ""
		for _, child := range sitemapModel.Childs {
			childCode += "\n" + renderSitemapEntry(templ, child)
		}

		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, childCode, 1)

	} else {

		// no childs
		rootCode = strings.Replace(rootCode, templates.ChildTemplatePlaceholder, "", 1)

	}

	return rootCode
}
