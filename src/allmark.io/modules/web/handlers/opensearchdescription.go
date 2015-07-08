// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/viewmodel"
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

func OpenSearchDescription(headerWriter header.HeaderWriter,
	openSearchDescriptionOrchestrator *orchestrator.OpenSearchDescriptionOrchestrator,
	templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_XML)

		// get the template
		hostname := getBaseURLFromRequest(r)
		openSearchDescriptionTemplate, err := templateProvider.GetSubTemplate(hostname, templates.OpenSearchDescriptionTemplateName)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		descriptionModel := openSearchDescriptionOrchestrator.GetDescriptionModel(hostname)
		openSearchDescription := getRenderedTemplateText(openSearchDescriptionTemplate, descriptionModel)

		fmt.Fprintf(w, "%s", openSearchDescription)
	})
}

func getRenderedTemplateText(templ *template.Template, model viewmodel.OpenSearchDescription) string {
	buffer := new(bytes.Buffer)
	renderTemplate(templ, model, buffer)
	return buffer.String()
}
