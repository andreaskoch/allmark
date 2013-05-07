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
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/view"
)

type Mapper func(parsedItem *parser.ParsedItem) view.Model

func Map(item *repository.Item, pathProvider *path.Provider) view.Model {

	// convert the item
	parsedItem, err := converter.Convert(item)
	if err != nil {
		return view.Error(fmt.Sprintf("%s", err), pathProvider.GetWebRoute(item))
	}

	switch itemType := parsedItem.MetaData.ItemType; itemType {
	case types.DocumentItemType:
		item.ViewModel = createDocumentMapperFunc(parsedItem, pathProvider)
		return item.ViewModel

	case types.MessageItemType:
		item.ViewModel = createMessageMapperFunc(parsedItem, pathProvider)
		return item.ViewModel

	case types.RepositoryItemType, types.CollectionItemType:
		item.ViewModel = createCollectionMapperFunc(parsedItem, pathProvider)
		return item.ViewModel

	default:
		item.ViewModel = view.Error(fmt.Sprintf("There is no mapper available for items of type %q", itemType), pathProvider.GetWebRoute(item))
		return item.ViewModel
	}

	panic("Unreachable")

}
