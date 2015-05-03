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
	"text/template"
)

type OpenSearchDescription struct {
	logger                            logger.Logger
	headerWriter                      header.HeaderWriter
	openSearchDescriptionOrchestrator *orchestrator.OpenSearchDescriptionOrchestrator
	templateProvider                  templates.Provider
}

func (handler *OpenSearchDescription) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_XML)

		// get the template
		hostname := getHostnameFromRequest(r)
		openSearchDescriptionTemplate, err := handler.templateProvider.GetSubTemplate(hostname, templates.OpenSearchDescriptionTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		descriptionModel := handler.openSearchDescriptionOrchestrator.GetDescriptionModel(hostname)
		openSearchDescription := getRenderedTemplateText(openSearchDescriptionTemplate, descriptionModel)

		fmt.Fprintf(w, "%s", openSearchDescription)
	}
}

func getRenderedTemplateText(templ *template.Template, model viewmodel.OpenSearchDescription) string {
	buffer := new(bytes.Buffer)
	renderTemplate(model, templ, buffer)
	return buffer.String()
}
