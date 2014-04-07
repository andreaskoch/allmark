// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeaheadorchestrator

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel/typeaheadviewmodel"
	"strings"
)

func NewSearchOrchestrator(searcher *search.ItemSearch, pathProvider paths.Pather) SearchOrchestrator {
	return SearchOrchestrator{
		searcher:     searcher,
		pathProvider: pathProvider,
	}
}

type SearchOrchestrator struct {
	searcher     *search.ItemSearch
	pathProvider paths.Pather
}

func (orchestrator *SearchOrchestrator) GetSuggestions(keywords string) []typeaheadviewmodel.SearchResult {

	// collect the search results
	typeAheadResults := make([]typeaheadviewmodel.SearchResult, 0)

	maximumNumberOfResults := 5

	if strings.TrimSpace(keywords) != "" {

		// execute the search
		searchResultItems := orchestrator.searcher.Search(keywords, maximumNumberOfResults)

		// prepare the result models
		for _, searchResult := range searchResultItems {
			typeAheadResults = append(typeAheadResults, orchestrator.createTypeAheadResultModel(searchResult))
		}

	}

	return typeAheadResults
}

func (orchestrator *SearchOrchestrator) createTypeAheadResultModel(searchResult search.SearchResult) typeaheadviewmodel.SearchResult {

	item := searchResult.Item

	// item location
	location := orchestrator.pathProvider.Path(item.Route().Value())

	return typeaheadviewmodel.SearchResult{
		Index: searchResult.Number,

		Title:       item.Title,
		Description: item.Description,
		Route:       location,
		Path:        item.Route().PrettyValue(),

		Value:  item.Title,
		Tokens: strings.Split(searchResult.StoreValue, " "),
	}
}
