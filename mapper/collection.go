// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func createCollectionMapperFunc(pathProvider *path.Provider, targetFormat string) Mapper {
	return func(item *repository.Item) view.Model {

		parsed, err := converter.Convert(item, targetFormat)
		if err != nil {
			return view.Error(fmt.Sprintf("%s", err))
		}

		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       parsed.Title,
			Description: parsed.Description,
			Content:     parsed.ConvertedContent,
			Entries:     getEntries(item, targetFormat),
			LanguageTag: getTwoLetterLanguageCode(parsed.MetaData.Language),
		}
	}
}

func getEntries(item *repository.Item, targetFormat string) []view.Model {

	entries := make([]view.Model, 0)

	// a path provider not relative to the repository but to the parent item
	relativePathProvider := path.NewProvider(item.Directory())

	for _, child := range item.Childs() {
		childMapper := New(child.Type, relativePathProvider, targetFormat)
		viewModel := childMapper(child)
		entries = append(entries, viewModel)
	}

	return entries
}
