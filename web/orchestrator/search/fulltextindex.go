// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/andreaskoch/allmark2/common/cleanup"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/bradleypeabody/fulltext"
	"strings"
)

type indexValueProvider func(item *model.Item) []string

func newIndex(logger logger.Logger, items []*model.Item, name string, indexValueFunc indexValueProvider) *FullTextIndex {

	index := &FullTextIndex{
		logger: logger,

		filepath:      fsutil.GetTempFileName(name),
		tempDirectory: fsutil.GetTempDirectory(),

		indexValueFunc: indexValueFunc,
	}

	go index.initialize(items)

	return index
}

type FullTextIndex struct {
	logger logger.Logger

	filepath      string
	tempDirectory string

	indexValueFunc indexValueProvider
}

func (index *FullTextIndex) Destroy() {

	// remove the index file
	cleanup.Now(index.filepath)

	// remove the temp directory
	cleanup.Now(index.tempDirectory)

	// self-destruct
	index = nil
}

func (index *FullTextIndex) Search(keywords string, maxiumNumberOfResults int) []Result {

	searcher, err := fulltext.NewSearcher(index.filepath)
	if err != nil {
		index.logger.Error(err.Error())
		return []Result{}
	}

	defer searcher.Close()

	searchResult, err := searcher.SimpleSearch(keywords, maxiumNumberOfResults)
	if err != nil {
		index.logger.Error(err.Error())
	}

	searchResults := make([]Result, 0)

	for number, v := range searchResult.Items {

		route, err := route.NewFromRequest(string(v.Id))
		if err != nil {
			index.logger.Warn("%s", err)
			continue
		}

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

func (index *FullTextIndex) initialize(items []*model.Item) {

	// fulltext search
	idx, err := fulltext.NewIndexer(index.tempDirectory)
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	for _, item := range items {

		doc := fulltext.IndexDoc{
			Id:         []byte(item.Route().Value()),              // unique identifier (the path to a webpage works...)
			StoreValue: []byte(item.Content),                      // bytes you want to be able to retrieve from search results
			IndexValue: getIndexValue(index.indexValueFunc(item)), // bytes you want to be split into words and indexed
		}

		idx.AddDoc(doc)
	}

	// when done, write out to final index
	f, err := fsutil.OpenFile(index.filepath)
	if err != nil {
		index.logger.Error(err.Error())
	}

	defer f.Close()

	err = idx.FinalizeAndWrite(f)
	if err != nil {
		index.logger.Error(err.Error())
	}
}

func getIndexValue(values []string) []byte {
	return []byte(strings.Join(values, " "))
}
