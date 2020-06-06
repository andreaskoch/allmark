// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/model"
	"github.com/andreaskoch/fulltext"
	"github.com/spf13/afero"
	"strings"
)

type indexValueProvider func(item *model.Item) []string

// newIndex creates a new FullTextIndex and initializes it with the given number of items.
func newIndex(logger logger.Logger, items []*model.Item, name string, indexValueFunc indexValueProvider) *FullTextIndex {

	index := &FullTextIndex{
		logger:         logger,
		filesystem:     &afero.MemMapFs{},
		indexValueFunc: indexValueFunc,
	}

	index.initialize(items)

	return index
}

// FullTextIndex indexes a given set if repository items and enables a full-text search on these items.
type FullTextIndex struct {
	logger logger.Logger

	filesystem afero.Fs

	indexValueFunc indexValueProvider
}

// Search scans the fulltext index for the given keywords and returns any matching search results.
func (index *FullTextIndex) Search(keywords string, maxiumNumberOfResults int) []Result {

	// open the search index file
	searchIndexFile, err := index.filesystem.Open("searchindex")
	if err != nil {
		index.logger.Error(err.Error())
		return []Result{}
	}

	// create a new searcher for the given search index
	searcher, err := fulltext.NewSearcher(searchIndexFile)
	if err != nil {
		index.logger.Error(err.Error())
		return []Result{}
	}

	defer searcher.Close()

	// peform the search
	searchResult, err := searcher.SimpleSearch(keywords, maxiumNumberOfResults)
	if err != nil {
		index.logger.Error(err.Error())
	}

	var searchResults []Result
	for number, v := range searchResult.Items {

		route := route.NewFromRequest(string(v.Id))

		// append the search results
		searchResults = append(searchResults, Result{
			Number: number + 1,

			Score:      v.Score,
			StoreValue: string(v.StoreValue),
			Route:      route,
		})

	}

	return searchResults
}

// initialize creates a fulltext index from the given repository items.
func (index *FullTextIndex) initialize(items []*model.Item) {

	// fulltext search
	indexer, err := fulltext.NewIndexer()
	if err != nil {
		index.logger.Error(err.Error())
		return
	}

	defer indexer.Close()

	for _, item := range items {

		doc := fulltext.IndexDoc{
			Id:         []byte(item.Route().Value()),              // unique identifier (the path to a webpage works...)
			StoreValue: []byte(item.Content),                      // bytes you want to be able to retrieve from search results
			IndexValue: getIndexValue(index.indexValueFunc(item)), // bytes you want to be split into words and indexed
		}

		indexer.AddDoc(doc)
	}

	// create a search index file
	searchIndexFile, err := index.filesystem.Create("searchindex")
	if err != nil {
		index.logger.Error(err.Error())
		return
	}

	// save the index to file
	if err := indexer.FinalizeAndWrite(searchIndexFile); err != nil {
		index.logger.Error(err.Error())
		return
	}

	defer searchIndexFile.Close()
}

func getIndexValue(values []string) []byte {
	return []byte(strings.Join(values, " "))
}
