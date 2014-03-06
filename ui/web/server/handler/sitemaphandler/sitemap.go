// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sitemaphandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	// "github.com/andreaskoch/allmark2/common/route"
	"bytes"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"net/http"
	"strings"
	"text/template"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) *SitemapHandler {

	templateProvider := templates.NewProvider(".")
	sitemapOrchestrator := orchestrator.NewSitemapOrchestrator(itemIndex)

	return &SitemapHandler{
		logger:              logger,
		itemIndex:           itemIndex,
		config:              config,
		patherFactory:       patherFactory,
		templateProvider:    templateProvider,
		sitemapOrchestrator: sitemapOrchestrator,
	}
}

type SitemapHandler struct {
	logger              logger.Logger
	itemIndex           *index.ItemIndex
	config              *config.Config
	patherFactory       paths.PatherFactory
	templateProvider    *templates.Provider
	sitemapOrchestrator orchestrator.SitemapOrchestrator
}

func (handler *SitemapHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the sitemap content template
		sitemapContentTemplate, err := handler.templateProvider.GetSubTemplate(templates.SitemapContentTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Content template not found. Error: %s", err)
			return
		}

		// get the sitemap template
		sitemapTemplate, err := handler.templateProvider.GetFullTemplate(templates.SitemapTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// render the sitemap content
		pathProvider := handler.patherFactory.Absolute("")
		sitemapContentModel := handler.sitemapOrchestrator.GetSitemap(pathProvider)
		sitemapContent := renderSitemapEntry(sitemapContentTemplate, &sitemapContentModel)

		sitemapPageModel := viewmodel.Model{
			Content: sitemapContent,
		}

		sitemapPageModel.Title = "Sitemap"
		sitemapPageModel.Description = "A list of all items in this repository."
		// sitemapPageModel.ToplevelNavigation = rootItem.ToplevelNavigation
		// sitemapPageModel.BreadcrumbNavigation = rootItem.BreadcrumbNavigation
		sitemapPageModel.Type = "sitemap"

		handlerutil.RenderTemplate(sitemapPageModel, sitemapTemplate, w)
	}
}

func renderSitemapEntry(templ *template.Template, sitemapModel *viewmodel.Sitemap) string {

	// render
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(sitemapModel, templ, buffer)

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
