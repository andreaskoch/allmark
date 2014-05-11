// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmlsitemaphandler

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
	"strings"
	"text/template"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory) *XmlSitemapHandler {

	templateProvider := templates.NewProvider(config.TemplatesFolder())
	xmlSitemapOrchestrator := orchestrator.NewXmlSitemapOrchestrator(itemIndex)

	return &XmlSitemapHandler{
		logger:                 logger,
		itemIndex:              itemIndex,
		config:                 config,
		patherFactory:          patherFactory,
		templateProvider:       templateProvider,
		xmlSitemapOrchestrator: xmlSitemapOrchestrator,
	}
}

type XmlSitemapHandler struct {
	logger                 logger.Logger
	itemIndex              *index.Index
	config                 *config.Config
	patherFactory          paths.PatherFactory
	templateProvider       *templates.Provider
	xmlSitemapOrchestrator orchestrator.XmlSitemapOrchestrator
}

func (handler *XmlSitemapHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the sitemap template
		xmlSitemapTemplate, err := handler.templateProvider.GetSubTemplate(templates.XmlSitemapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		sitemapWrapper := renderSitemapWrapper(xmlSitemapTemplate)

		// get the sitemap content template
		xmlSitemapContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.XmlSitemapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		// render the sitemap content
		entries := handler.xmlSitemapOrchestrator.GetSitemapEntires(pathProvider)

		sitemapContent := renderSitemapEntries(xmlSitemapContentTemplate, entries)

		sitemapWrapper = strings.Replace(sitemapWrapper, templates.ChildTemplatePlaceholder, sitemapContent, 1)

		fmt.Fprintf(w, "%s", sitemapWrapper)
	}
}

func renderSitemapWrapper(templ *template.Template) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(nil, templ, buffer)
	return buffer.String()
}

func renderSitemapEntries(templ *template.Template, sitemapEntries []viewmodel.XmlSitemapEntry) string {

	rootCode := ""
	for _, entry := range sitemapEntries {
		buffer := new(bytes.Buffer)
		handlerutil.RenderTemplate(entry, templ, buffer)
		rootCode += "\n" + buffer.String()
	}

	return rootCode
}
