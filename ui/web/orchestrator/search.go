// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

var (
	itemsPerPage = 5
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
		panic("Invalid page number")
	}

	// determine start item
	startItemNumber := itemsPerPage * (page - 1)

	// determine end item
	endItemNumber := itemsPerPage * page

	// collect the search results
	searchResults := make([]viewmodel.SearchResult, 0)
	for _, searchResult := range orchestrator.fullTextIndex.Search(keywords) {

		// paging
		currentNumberOfItems := len(searchResults)
		if currentNumberOfItems < startItemNumber || currentNumberOfItems >= endItemNumber {
			continue
		}

		searchResults = append(searchResults, orchestrator.createSearchResultModel(searchResult))
	}

	return viewmodel.Search{
		Query:        keywords,
		Page:         page,
		ItemsPerPage: itemsPerPage,
		Results:      searchResults,
	}
}

func (orchestrator *SearchOrchestrator) createSearchResultModel(searchResult search.SearchResult) viewmodel.SearchResult {

	item := searchResult.Item

	// item location
	location := orchestrator.pathProvider.Path(item.Route().Value())

	// last modified date
	lastModifiedDate := ""
	if item.MetaData != nil && item.MetaData.LastModifiedDate != nil {
		lastModifiedDate = item.MetaData.LastModifiedDate.Format("2006-01-02")
	}

	return viewmodel.SearchResult{
		Title:       item.Title,
		Description: searchResult.StoreValue,
		Route:       location,
		PubDate:     lastModifiedDate,
	}
}
