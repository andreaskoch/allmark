// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"github.com/elWyatt/allmark/web/view/templates"
	"fmt"
	"net/http"
)

// OpenSearchDescription returns a opensearch description http handler.
func OpenSearchDescription(headerWriter header.HeaderWriter,
	openSearchDescriptionOrchestrator *orchestrator.OpenSearchDescriptionOrchestrator,
	templateProvider templates.Provider) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_XML)

		// get the template
		hostname := getBaseURLFromRequest(r)
		openSearchDescriptionTemplate, err := templateProvider.GetOpenSearchDescriptionTemplate(hostname)
		if err != nil {
			fmt.Fprintf(w, "Template not found. Error: %s", err)
			return
		}

		descriptionModel := openSearchDescriptionOrchestrator.GetDescriptionModel(hostname)
		renderTemplate(openSearchDescriptionTemplate, descriptionModel, w)
	})
}
