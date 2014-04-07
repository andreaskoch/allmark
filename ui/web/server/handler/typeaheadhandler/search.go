// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeaheadhandler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel/typeaheadviewmodel"
	"io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex, searcher *search.ItemSearch) *TypeAheadHandler {

	// search
	searchResultPathProvider := patherFactory.Absolute("/")
	typeAheadOrchestrator := orchestrator.NewTypeAheadOrchestrator(searcher, searchResultPathProvider)

	return &TypeAheadHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		config:                config,
		patherFactory:         patherFactory,
		typeAheadOrchestrator: typeAheadOrchestrator,
	}
}

type TypeAheadHandler struct {
	logger                logger.Logger
	itemIndex             *index.ItemIndex
	config                *config.Config
	patherFactory         paths.PatherFactory
	typeAheadOrchestrator orchestrator.TypeAheadOrchestrator
}

func (handler *TypeAheadHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the query parameter
		query, _ := handlerutil.GetQueryParameterFromUrl(*r.URL)

		// get the suggestions
		typeAheadResults := handler.typeAheadOrchestrator.GetSuggestions(query)

		// set content type to json
		w.Header().Set("Content-Type", "application/json")

		// convert to json
		writeTypeAheadResults(w, typeAheadResults)
	}
}

func writeTypeAheadResults(writer io.Writer, typeAheadResults []typeaheadviewmodel.SearchResult) error {
	bytes, err := json.MarshalIndent(typeAheadResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
