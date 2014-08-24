// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"strings"
)

type TypeAheadOrchestrator struct {
	*Orchestrator
}

func (orchestrator *TypeAheadOrchestrator) GetSuggestions(keywords string) []viewmodel.TypeAhead {

	// collect the search results
	typeAheadResults := make([]viewmodel.TypeAhead, 0)

	maximumNumberOfResults := 5

	if strings.TrimSpace(keywords) != "" {

		// execute the search
		searchResultItems := orchestrator.repository.Search(keywords, maximumNumberOfResults)

		// prepare the result models
		for _, searchResult := range searchResultItems {
			typeAheadResults = append(typeAheadResults, orchestrator.createTypeAheadResultModel(searchResult))
		}

	}

	return typeAheadResults
}

func (orchestrator *TypeAheadOrchestrator) createTypeAheadResultModel(searchResult dataaccess.SearchResult) viewmodel.TypeAhead {

	item := orchestrator.parseItem(searchResult.Item)
	if item == nil {
		return viewmodel.TypeAhead{}
	}

	// item location
	location := orchestrator.itemPather().Path(item.Route().Value())

	return viewmodel.TypeAhead{
		Index: searchResult.Number,

		Title:       item.Title,
		Description: item.Description,
		Route:       location,
		Path:        item.Route().OriginalValue(),

		Value:  item.Title,
		Tokens: strings.Split(searchResult.StoreValue, " "),
	}
}
