// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"io"
	"net/http"
)

type TypeAhead struct {
	logger logger.Logger

	typeAheadOrchestrator orchestrator.TypeAheadOrchestrator
}

func (handler *TypeAhead) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the query parameter
		query, _ := getQueryParameterFromUrl(*r.URL)

		// get the suggestions
		searchResults := handler.typeAheadOrchestrator.GetSuggestions(query)

		// set content type to json
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/json")

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
