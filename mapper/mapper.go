// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

type Mapper func(parsedItem *parser.Result) view.Model

func Map(item *repository.Item, pathProvider *path.Provider, targetFormat string) view.Model {

	// convert the item
	parsedItem, err := converter.Convert(item, targetFormat)
	if err != nil {
		return view.Error(fmt.Sprintf("%s", err), pathProvider.GetWebRoute(item))
	}

	switch parsedItem.MetaData.ItemType {
	case parser.DocumentItemType:
		return createDocumentMapperFunc(parsedItem, pathProvider, targetFormat)

	case parser.MessageItemType:
		return createMessageMapperFunc(parsedItem, pathProvider, targetFormat)

	case parser.RepositoryItemType, parser.CollectionItemType:
		return createCollectionMapperFunc(parsedItem, pathProvider, targetFormat)
	}

	return view.Error(fmt.Sprintf("There is no mapper available for items of type %q", parsedItem.MetaData.ItemType), pathProvider.GetWebRoute(item))
}
