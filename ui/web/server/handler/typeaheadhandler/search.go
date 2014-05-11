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
	"github.com/andreaskoch/allmark2/ui/web/orchestrator/typeaheadorchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel/typeaheadviewmodel"
	"io"
	"net/http"
)

func NewSearchHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.Index, searcher *search.ItemSearch) *SearchHandler {

	// search
	searchResultPathProvider := patherFactory.Absolute("/")
	searchOrchestrator := typeaheadorchestrator.NewSearchOrchestrator(searcher, searchResultPathProvider)

	return &SearchHandler{
		logger:             logger,
		itemIndex:          itemIndex,
		config:             config,
		patherFactory:      patherFactory,
		searchOrchestrator: &searchOrchestrator,
	}
}

type SearchHandler struct {
	logger             logger.Logger
	itemIndex          *index.Index
	config             *config.Config
	patherFactory      paths.PatherFactory
	searchOrchestrator *typeaheadorchestrator.SearchOrchestrator
}

func (handler *SearchHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the query parameter
		query, _ := handlerutil.GetQueryParameterFromUrl(*r.URL)

		// get the suggestions
		searchResults := handler.searchOrchestrator.GetSuggestions(query)

		// set content type to json
		w.Header().Set("Content-Type", "application/json")

		// convert to json
		writeSearchResults(w, searchResults)
	}
}

func writeSearchResults(writer io.Writer, searchResults []typeaheadviewmodel.SearchResult) error {
	bytes, err := json.MarshalIndent(searchResults, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
