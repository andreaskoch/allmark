// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/model"
	"strings"
)

func NewItemSearch(logger logger.Logger, itemIndex *index.Index) *ItemSearch {

	return &ItemSearch{
		logger:    logger,
		itemIndex: itemIndex,

		routesFullTextIndex:      newIndex(logger, itemIndex, "route", routes),
		tagsFullTextIndex:        newIndex(logger, itemIndex, "tags", tags),
		itemContentFullTextIndex: newIndex(logger, itemIndex, "content", content),
	}
}

func routes(item *model.Item) []string {
	if item == nil {
		return []string{}
	}

	return item.Route().Components()
}

func tags(item *model.Item) []string {
	if item == nil || item.MetaData == nil {
		return []string{}
	}

	keywords := make([]string, 0)
	for _, tag := range item.MetaData.Tags {
		keywords = append(keywords, tag.Name())
	}

	return keywords
}

func content(item *model.Item) []string {

	if item == nil {
		return []string{}
	}

	keywords := make([]string, 0)
	keywords = append(keywords, item.Title)
	keywords = append(keywords, item.Description)
	keywords = append(keywords, item.Route().Components()...)
	keywords = append(keywords, item.Content)
	return keywords
}

type ItemSearch struct {
	logger    logger.Logger
	itemIndex *index.Index

	routesFullTextIndex      *FullTextIndex
	tagsFullTextIndex        *FullTextIndex
	itemContentFullTextIndex *FullTextIndex
}

func (itemSearch *ItemSearch) Search(keywords string, maxiumNumberOfResults int) []SearchResult {

	// routes
	if isRouteSearch(keywords) {
		routeComponents := strings.Replace(keywords, "/", " ", -1)
		return itemSearch.routesFullTextIndex.Search(routeComponents, maxiumNumberOfResults)
	}

	// tags
	if isTagSearch(keywords) {
		tagWithoutPrefix := strings.Replace(keywords, "#", "", -1)
		return itemSearch.tagsFullTextIndex.Search(tagWithoutPrefix, maxiumNumberOfResults)
	}

	// items
	return itemSearch.itemContentFullTextIndex.Search(keywords, maxiumNumberOfResults)
}

func (itemSearch *ItemSearch) Update() {
	go itemSearch.routesFullTextIndex.Update()
	go itemSearch.tagsFullTextIndex.Update()
	go itemSearch.itemContentFullTextIndex.Update()
}

func isRouteSearch(keywords string) bool {
	return strings.HasPrefix(keywords, "/")
}

func isTagSearch(keywords string) bool {
	return strings.HasPrefix(keywords, "#")
}
