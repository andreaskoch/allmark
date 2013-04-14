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

func createCollectionMapperFunc(pathProvider *path.Provider, converterFactory func(item *repository.Item) func() string) func(item *repository.Item) view.Model {
	return func(item *repository.Item) view.Model {

		converter := converterFactory(item)
		html := converter()

		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       item.Title,
			Description: item.Description,
			Content:     html,
			Entries:     getEntries(pathProvider, converterFactory, item),
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}
}

func getEntries(pathProvider *path.Provider, converterFactory func(item *repository.Item) func() string, item *repository.Item) []view.Model {

	entries := make([]view.Model, 0)

	for _, child := range item.Childs() {
		if mapperFunc, err := GetMapper(pathProvider, converterFactory, child.Type); err == nil {
			viewModel := mapperFunc(child)
			entries = append(entries, viewModel)
		} else {
			fmt.Println(err)
		}

	}

	return entries
}
