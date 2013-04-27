// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/view"
)

func createCollectionMapperFunc(parsedItem *parser.ParsedItem, pathProvider *path.Provider, targetFormat string) view.Model {

	return view.Model{
		Path:        pathProvider.GetWebRoute(parsedItem),
		Title:       parsedItem.Title,
		Description: parsedItem.Description,
		Content:     parsedItem.ConvertedContent,
		Entries:     getEntries(parsedItem, pathProvider.UseTempDir(), targetFormat),
		Type:        parsedItem.MetaData.ItemType,
		LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
	}

}

func getEntries(parsedItem *parser.ParsedItem, useTempDir bool, targetFormat string) []view.Model {

	viewModels := make([]view.Model, 0)

	// a path provider not relative to the repository but to the parent item
	relativePathProvider := path.NewProvider(parsedItem.Directory(), useTempDir)

	for _, child := range parsedItem.Childs() {
		viewModel := Map(child, relativePathProvider, targetFormat)
		viewModels = append(viewModels, viewModel)
	}

	return viewModels
}
