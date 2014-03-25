// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"strings"
)

var (
	itemsPerPage = 10
)

func NewSearchOrchestrator(fullTextIndex *search.FullTextIndex, pathProvider paths.Pather) SearchOrchestrator {
	return SearchOrchestrator{
		fullTextIndex: fullTextIndex,
		pathProvider:  pathProvider,
	}
}

type SearchOrchestrator struct {
	fullTextIndex *search.FullTextIndex
	pathProvider  paths.Pather
}

func (orchestrator *SearchOrchestrator) GetSearchResults(keywords string, page int) viewmodel.Search {

	// validate page number
	if page < 1 {
		panic("Invalid page number.")
	}

	// determine start item
	startItemNumber := itemsPerPage * (page - 1)

	// determine end item
	endItemNumber := itemsPerPage * page

	// collect the search results
	searchResults := make([]viewmodel.SearchResult, 0)

	maximumNumberOfResults := 100
	totalResultCount := 0

	if strings.TrimSpace(keywords) != "" {

		// execute the search
		searchResultItems := orchestrator.fullTextIndex.Search(keywords, maximumNumberOfResults)

		// count the number of search results
		totalResultCount = len(searchResultItems)

		// prepare the result models
		for currentNumberOfItems, searchResult := range searchResultItems {

			// paging
			if currentNumberOfItems < startItemNumber || currentNumberOfItems >= endItemNumber {
				continue
			}

			searchResults = append(searchResults, orchestrator.createSearchResultModel(searchResult))
		}

	}

	return viewmodel.Search{
		Query:   keywords,
		Results: searchResults,

		Page:         page,
		ItemsPerPage: itemsPerPage,

		ResultCount:      len(searchResults),
		TotalResultCount: totalResultCount,
	}
}

func (orchestrator *SearchOrchestrator) createSearchResultModel(searchResult search.SearchResult) viewmodel.SearchResult {

	item := searchResult.Item

	// item location
	location := orchestrator.pathProvider.Path(item.Route().Value())

	return viewmodel.SearchResult{
		Title:       item.Title,
		Description: item.Description,
		Route:       location,
		Path:        item.Route().PrettyValue(),
	}
}
