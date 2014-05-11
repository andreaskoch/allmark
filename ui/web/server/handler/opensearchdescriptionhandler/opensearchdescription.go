// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opensearchdescriptionhandler

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
	"text/template"
)

func New(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.Index) *OpenSearchDescriptionHandler {

	templateProvider := templates.NewProvider(config.TemplatesFolder())
	openSearchDescriptionOrchestrator := orchestrator.NewOpenSearchDescriptionOrchestrator(itemIndex)

	return &OpenSearchDescriptionHandler{
		logger:                            logger,
		itemIndex:                         itemIndex,
		config:                            config,
		patherFactory:                     patherFactory,
		templateProvider:                  templateProvider,
		openSearchDescriptionOrchestrator: openSearchDescriptionOrchestrator,
	}
}

type OpenSearchDescriptionHandler struct {
	logger                            logger.Logger
	itemIndex                         *index.Index
	config                            *config.Config
	patherFactory                     paths.PatherFactory
	templateProvider                  *templates.Provider
	openSearchDescriptionOrchestrator orchestrator.OpenSearchDescriptionOrchestrator
}

func (handler *OpenSearchDescriptionHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the template
		openSearchDescriptionTemplate, err := handler.templateProvider.GetSubTemplate(templates.OpenSearchDescriptionTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		// prepare a path provider which includes the hostname
		hostname := handlerutil.GetHostnameFromRequest(r)
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := handler.patherFactory.Absolute(addressPrefix)

		descriptionModel := handler.openSearchDescriptionOrchestrator.GetDescriptionModel(pathProvider)
		openSearchDescription := renderTemplate(openSearchDescriptionTemplate, &descriptionModel)

		fmt.Fprintf(w, "%s", openSearchDescription)
	}
}

func renderTemplate(templ *template.Template, model *viewmodel.OpenSearchDescription) string {
	buffer := new(bytes.Buffer)
	handlerutil.RenderTemplate(model, templ, buffer)
	return buffer.String()
}
