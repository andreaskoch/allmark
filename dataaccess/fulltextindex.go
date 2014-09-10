// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/bradleypeabody/fulltext"
	"io/ioutil"
	"os"
	"strings"
)

type indexValueProvider func(item *Item) []string

func newIndex(logger logger.Logger, repository Repository, name string, indexValueFunc indexValueProvider) *FullTextIndex {

	index := &FullTextIndex{
		logger: logger,

		repository: repository,

		filepath:       getTempFileName(name),
		indexValueFunc: indexValueFunc,
	}

	go index.initialize()

	return index
}

type FullTextIndex struct {
	logger logger.Logger

	repository Repository

	filepath       string
	indexValueFunc indexValueProvider
}

func (index *FullTextIndex) Search(keywords string, maxiumNumberOfResults int) []SearchResult {

	searcher, err := fulltext.NewSearcher(index.filepath)
	if err != nil {
		index.logger.Error(err.Error())
	}

	defer searcher.Close()

	searchResult, err := searcher.SimpleSearch(keywords, maxiumNumberOfResults)
	if err != nil {
		index.logger.Error(err.Error())
	}

	searchResults := make([]SearchResult, 0)

	for number, v := range searchResult.Items {

		route, err := route.NewFromRequest(string(v.Id))
		if err != nil {
			index.logger.Warn("%s", err)
			continue
		}

		item, exists := index.repository.Item(route)

		// skip if the item was not found in the repository
		if !exists {
			continue
		}

		// append the search results
		searchResults = append(searchResults, SearchResult{
			Number: number + 1,

			Score:      v.Score,
			StoreValue: string(v.StoreValue),
			Item:       item,
		})

	}

	return searchResults
}

func (index *FullTextIndex) initialize() {

	// fulltext search
	idx, err := fulltext.NewIndexer("")
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	for _, item := range index.repository.Items() {

		doc := fulltext.IndexDoc{
			Id:         []byte(item.Route().Value()),              // unique identifier (the path to a webpage works...)
			StoreValue: []byte(getContentFromItem(item)),          // bytes you want to be able to retrieve from search results
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
