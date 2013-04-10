// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func createCollectionMapperFunc(pathProvider *path.Provider, converterFunc func(item *repository.Item) string) func(item *repository.Item) view.Model {
	return func(item *repository.Item) view.Model {

		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       item.Title,
			Description: item.Description,
			Content:     converterFunc(item),
			Entries:     getEntries(pathProvider, converterFunc, item),
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}
}

func getEntries(pathProvider *path.Provider, converterFunc func(item *repository.Item) string, item *repository.Item) []view.Model {

	entries := make([]view.Model, 0)

	for _, child := range item.ChildItems {
		if mapperFunc, err := GetMapper(pathProvider, converterFunc, child); err == nil {
			viewModel := mapperFunc(child)
			entries = append(entries, viewModel)
		} else {
			fmt.Println(err)
		}

	}

	return entries
}
