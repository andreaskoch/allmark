// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/bradleypeabody/fulltext"
)

func NewIndex(logger logger.Logger, itemIndex *index.ItemIndex) *FullTextIndex {
	return &FullTextIndex{
		logger:    logger,
		itemIndex: itemIndex,
		filepath:  "index",
	}
}

type FullTextIndex struct {
	logger    logger.Logger
	itemIndex *index.ItemIndex
	filepath  string
}

func (index *FullTextIndex) Update() {

	// fulltext search
	idx, err := fulltext.NewIndexer("")
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	for _, item := range index.itemIndex.Items() {

		doc := fulltext.IndexDoc{
			Id:         []byte(item.Route().Value()), // unique identifier (the path to a webpage works...)
			StoreValue: []byte(item.Title),           // bytes you want to be able to retrieve from search results
			IndexValue: getContent(item),             // bytes you want to be split into words and indexed
		}

		idx.AddDoc(doc)
	}

	// when done, write out to final index
	f, err := fsutil.OpenFile(index.filepath)
	if err != nil {
		panic(err)
	}

	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}
}

func getContent(item *model.Item) []byte {

	content := item.Title
	content += " " + item.Description
	content += " " + item.Content

	return []byte(content)
}

func (index *FullTextIndex) Search(keyword string) []SearchResult {

	searcher, err := fulltext.NewSearcher(index.filepath)
	if err != nil {
		panic(err)
	}

	defer searcher.Close()

	searchResult, err := searcher.SimpleSearch(keyword, 5)
	if err != nil {
		panic(err)
	}

	index.logger.Debug("%s", keyword)
	index.logger.Debug("%s", len(searchResult.Items))

	searchResults := make([]SearchResult, 0)

	for _, v := range searchResult.Items {

		route, err := route.NewFromRequest(string(v.Id))
		if err != nil {
			index.logger.Warn("%s", err)
			continue
		}

		if item, isMatch := index.itemIndex.IsMatch(*route); isMatch {
			searchResults = append(searchResults, SearchResult{
				Score:      v.Score,
				StoreValue: string(v.StoreValue),
				Item:       item,
			})
		}

	}

	return searchResults
}
