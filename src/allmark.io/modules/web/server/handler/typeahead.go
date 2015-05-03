// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
)

type TypeAhead struct {
	logger                logger.Logger
	headerWriter          header.HeaderWriter
	typeAheadOrchestrator *orchestrator.TypeAheadOrchestrator
}

func (handler *TypeAhead) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		handler.headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// get the suggestions
		query, _ := getQueryParameterFromUrl(*r.URL)
		searchResults := handler.typeAheadOrchestrator.GetSuggestions(query)

		// convert to json
		writeSearchResults(w, searchResults)
	}
}

func writeSearchResults(writer io.Writer, searchResults []viewmodel.TypeAhead) error {
	bytes, err := json.MarshalIndent(searchResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
