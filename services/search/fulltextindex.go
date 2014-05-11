// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/bradleypeabody/fulltext"
	"io/ioutil"
	"os"
	"strings"
)

type indexValueProvider func(item *model.Item) []string

func newIndex(logger logger.Logger, itemIndex *index.Index, name string, indexValueFunc indexValueProvider) *FullTextIndex {

	return &FullTextIndex{
		logger:         logger,
		itemIndex:      itemIndex,
		filepath:       getTempFileName(name),
		indexValueFunc: indexValueFunc,
	}
}

type FullTextIndex struct {
	logger         logger.Logger
	itemIndex      *index.Index
	filepath       string
	indexValueFunc indexValueProvider
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
			Id:         []byte(item.Route().Value()),              // unique identifier (the path to a webpage works...)
			StoreValue: []byte(item.Title),                        // bytes you want to be able to retrieve from search results
			IndexValue: getIndexValue(index.indexValueFunc(item)), // bytes you want to be split into words and indexed
		}

		idx.AddDoc(doc)
	}

	// when done, write out to final index
	f, err := fsutil.OpenFile(index.filepath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}
}

func (index *FullTextIndex) Search(keywords string, maxiumNumberOfResults int) []SearchResult {

	searcher, err := fulltext.NewSearcher(index.filepath)
	if err != nil {
		panic(err)
	}

	defer searcher.Close()

	searchResult, err := searcher.SimpleSearch(keywords, maxiumNumberOfResults)
	if err != nil {
		panic(err)
	}

	searchResults := make([]SearchResult, 0)

	for number, v := range searchResult.Items {

		route, err := route.NewFromRequest(string(v.Id))
		if err != nil {
			index.logger.Warn("%s", err)
			continue
		}

		if item, isMatch := index.itemIndex.IsMatch(*route); isMatch {
			searchResults = append(searchResults, SearchResult{
				Number: number + 1,

				Score:      v.Score,
				StoreValue: string(v.StoreValue),
				Item:       item,
			})
		}

	}

	return searchResults
}

func getTempFileName(prefix string) string {
	file, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%s-index", prefix))
	if err != nil {
		panic(err)
	}

	defer file.Close()

	return file.Name()
}

func getIndexValue(values []string) []byte {
	return []byte(strings.Join(values, " "))
}
