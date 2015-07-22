// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"strings"
)

// Result is the model returned by the fulltext index's Search function.
type Result struct {
	Route route.Route

	Number     int
	Score      int64
	StoreValue string
}

// NewItemSearch creates a new repository item searcher.
func NewItemSearch(logger logger.Logger, items []*model.Item) *ItemSearch {

	return &ItemSearch{
		logger: logger,

		routesFullTextIndex:      newIndex(logger, items, "route", itemRouteKeywordProvider),
		itemContentFullTextIndex: newIndex(logger, items, "content", itemContentKeywordProvider),
	}
}

// itemRouteKeywordProvider returns a list of keywords for the fulltext index
// from the given items' route.
func itemRouteKeywordProvider(item *model.Item) []string {
	if item == nil {
		return []string{}
	}

	// item route components
	routeComponents := item.Route().Components()

	// file route components
	for _, file := range item.Files() {
		routeComponents = append(routeComponents, file.Route().Components()...)
	}

	return routeComponents
}

// itemContentKeywordProvider returns a list of keywords for the fulltext index
// from the given items' content.
func itemContentKeywordProvider(item *model.Item) []string {

	if item == nil {
		return []string{}
	}

	var keywords []string
	keywords = append(keywords, getContentFromItem(item))

	return keywords
}

// ItemSearch creates a fulltext index for given set of repository items and provides
// the ability to search over this index.
type ItemSearch struct {
	logger logger.Logger

	routesFullTextIndex      *FullTextIndex
	itemContentFullTextIndex *FullTextIndex
}

// Search returns a set of Result models that match specified keywords.
func (itemSearch *ItemSearch) Search(keywords string, maxiumNumberOfResults int) []Result {

	// routes
	if isRouteSearch(keywords) {
		routeComponents := strings.Replace(keywords, "/", " ", -1)
		return itemSearch.routesFullTextIndex.Search(routeComponents, maxiumNumberOfResults)
	}

	// items
	return itemSearch.itemContentFullTextIndex.Search(keywords, maxiumNumberOfResults)
}

// getContentFromItem returns the content from the given repository item.
func getContentFromItem(item *model.Item) string {

	return item.Description + " " + item.Content
}

// isRouteSearch detects based on the given keywords
// if the search targeted towards a route.
func isRouteSearch(keywords string) bool {
	return strings.HasPrefix(keywords, "/")
}
