// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"io"
	"net/http"
)

type TypeAhead struct {
	logger logger.Logger

	typeAheadOrchestrator *orchestrator.TypeAheadOrchestrator
}

func (handler *TypeAhead) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "application/json; charset=utf-8")
		header.NoCache(w, r)
		header.VaryAcceptEncoding(w, r)

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
