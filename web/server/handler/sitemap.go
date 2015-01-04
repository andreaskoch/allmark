// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"net/http"
	"strings"
	"text/template"
)

type Sitemap struct {
	logger logger.Logger

	navigationOrchestrator *orchestrator.NavigationOrchestrator
	sitemapOrchestrator    *orchestrator.SitemapOrchestrator

	templateProvider templates.Provider
}

func (self *Sitemap) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/html; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		hostname := getHostnameFromRequest(r)

		// get the sitemap content template
		sitemapContentTemplate, err := self.templateProvider.GetSubTemplate(hostname, templates.SitemapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// get the sitemap template
		sitemapTemplate, err := self.templateProvider.GetFullTemplate(hostname, templates.SitemapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		sitemapContentModel := self.sitemapOrchestrator.GetSitemap()
		sitemapContent := renderSitemapEntry(sitemapContentTemplate, sitemapContentModel)

		sitemapPageModel := viewmodel.Model{
			Content: sitemapContent,
		}

		sitemapPageModel.Type = "sitemap"
		sitemapPageModel.Title = "Sitemap"
		sitemapPageModel.Description = "A list of all items in this repository."
		sitemapPageModel.ToplevelNavigation = self.navigationOrchestrator.GetToplevelNavigation()
		sitemapPageModel.BreadcrumbNavigation = self.navigationOrchestrator.GetBreadcrumbNavigation(route.New())

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
