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

type XmlSitemap struct {
	logger logger.Logger

	xmlSitemapOrchestrator *orchestrator.XmlSitemapOrchestrator

	templateProvider templates.Provider
}

func (handler *XmlSitemap) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/xml; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		hostname := getHostnameFromRequest(r)

		// get the sitemap template
		xmlSitemapTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.XmlSitemapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		sitemapWrapper := renderSitemapWrapper(xmlSitemapTemplate)

		// get the sitemap content template
		xmlSitemapContentTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.XmlSitemapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		entries := handler.xmlSitemapOrchestrator.GetSitemapEntires(hostname)

		sitemapContent := renderSitemapEntries(xmlSitemapContentTemplate, entries)

		sitemapWrapper = strings.Replace(sitemapWrapper, templates.ChildTemplatePlaceholder, sitemapContent, 1)

		fmt.Fprintf(w, "%s", sitemapWrapper)
	}
}

func renderSitemapWrapper(templ *template.Template) string {
	buffer := new(bytes.Buffer)
	renderTemplate(nil, templ, buffer)
	return buffer.String()
}

func renderSitemapEntries(templ *template.Template, sitemapEntries []viewmodel.XmlSitemapEntry) string {

	rootCode := ""
	for _, entry := range sitemapEntries {
		buffer := new(bytes.Buffer)
		renderTemplate(entry, templ, buffer)
		rootCode += "\n" + buffer.String()
	}

	return rootCode
}
