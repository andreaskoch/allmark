// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/web/header"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
)

func TypeAhead(headerWriter header.HeaderWriter, typeAheadOrchestrator *orchestrator.TypeAheadOrchestrator) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// get the suggestions
		query, _ := getQueryParameterFromURL(*r.URL)
		searchResults := typeAheadOrchestrator.GetSuggestions(query)

		// convert to json
		writeSearchResults(w, searchResults)
	})

}

func writeSearchResults(writer io.Writer, searchResults []viewmodel.TypeAhead) error {
	bytes, err := json.MarshalIndent(searchResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
