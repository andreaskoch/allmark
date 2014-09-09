// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"bufio"
	"bytes"
	"github.com/andreaskoch/allmark2/common/logger"
	"io"
	"strings"
)

func NewItemSearch(logger logger.Logger, repository Repository) *ItemSearch {

	return &ItemSearch{
		logger: logger,

		repository: repository,

		routesFullTextIndex:      newIndex(logger, repository, "route", itemRouteKeywordProvider),
		itemContentFullTextIndex: newIndex(logger, repository, "content", itemContentKeywordProvider),
	}
}

func itemRouteKeywordProvider(item *Item) []string {
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

func itemContentKeywordProvider(item *Item) []string {

	if item == nil {
		return []string{}
	}

	keywords := make([]string, 0)
	keywords = append(keywords, getContentFromItem(item))

	return keywords
}

type ItemSearch struct {
	logger     logger.Logger
	repository Repository

	routesFullTextIndex      *FullTextIndex
	itemContentFullTextIndex *FullTextIndex
}

func (itemSearch *ItemSearch) Search(keywords string, maxiumNumberOfResults int) []SearchResult {

	// routes
	if isRouteSearch(keywords) {
		routeComponents := strings.Replace(keywords, "/", " ", -1)
		return itemSearch.routesFullTextIndex.Search(routeComponents, maxiumNumberOfResults)
	}

	// items
	return itemSearch.itemContentFullTextIndex.Search(keywords, maxiumNumberOfResults)
}

func getContentFromItem(item *Item) string {

	// fetch the item data
	byteBuffer := new(bytes.Buffer)
	dataWriter := bufio.NewWriter(byteBuffer)

	contentReader := func(content io.ReadSeeker) error {
		_, err := io.Copy(dataWriter, content)
		dataWriter.Flush()
		return err
	}

	if err := item.Data(contentReader); err != nil {
		return ""
	}

	return byteBuffer.String()
}

func isRouteSearch(keywords string) bool {
	return strings.HasPrefix(keywords, "/")
}
