// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/view"
)

func createDocumentMapperFunc(parsedItem *parser.ParsedItem, pathProvider *path.Provider, targetFormat string) view.Model {

	return view.Model{
		Path:        pathProvider.GetWebRoute(parsedItem),
		Title:       parsedItem.Title,
		Description: parsedItem.Description,
		Content:     parsedItem.ConvertedContent,
		LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
		Type:        parsedItem.MetaData.ItemType,
	}

}
