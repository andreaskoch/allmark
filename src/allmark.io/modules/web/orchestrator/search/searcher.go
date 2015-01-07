// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"strings"
)

type Result struct {
	Number int

	Score      int64
	StoreValue string
	Route      route.Route
}

func NewItemSearch(logger logger.Logger, items []*model.Item) *ItemSearch {

	return &ItemSearch{
		logger: logger,

		routesFullTextIndex:      newIndex(logger, items, "route", itemRouteKeywordProvider),
		itemContentFullTextIndex: newIndex(logger, items, "content", itemContentKeywordProvider),
	}
}

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

func itemContentKeywordProvider(item *model.Item) []string {

	if item == nil {
		return []string{}
	}

	keywords := make([]string, 0)
	keywords = append(keywords, getContentFromItem(item))

	return keywords
}

type ItemSearch struct {
	logger logger.Logger

	routesFullTextIndex      *FullTextIndex
	itemContentFullTextIndex *FullTextIndex
}

func (itemSearch *ItemSearch) Destroy() {

	// destroy the indizes
	itemSearch.itemContentFullTextIndex.Destroy()
	itemSearch.routesFullTextIndex.Destroy()

	// self-destruct
	itemSearch = nil
}

func (itemSearch *ItemSearch) Search(keywords string, maxiumNumberOfResults int) []Result {

	// routes
	if isRouteSearch(keywords) {
		routeComponents := strings.Replace(keywords, "/", " ", -1)
		return itemSearch.routesFullTextIndex.Search(routeComponents, maxiumNumberOfResults)
	}

	// items
	return itemSearch.itemContentFullTextIndex.Search(keywords, maxiumNumberOfResults)
}

func getContentFromItem(item *model.Item) string {

	return item.Description + " " + item.Content
}

func isRouteSearch(keywords string) bool {
	return strings.HasPrefix(keywords, "/")
}
