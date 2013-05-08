// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/view"
)

type Mapper func(parsedItem *parser.ParsedItem) view.Model

func Map(item *repository.Item) {

	// convert the item
	parsedItem, err := converter.Convert(item)
	if err != nil {
		item.Model = view.Error(fmt.Sprintf("%s", err), parsedItem.PathProvider().GetWebRoute(item))
		return
	}

	switch itemType := parsedItem.MetaData.ItemType; itemType {
	case types.DocumentItemType, types.RepositoryItemType, types.CollectionItemType:
		item.Model = createDocumentMapperFunc(parsedItem)
		return

	case types.MessageItemType:
		item.Model = createMessageMapperFunc(parsedItem)
		return

	default:
		item.Model = view.Error(fmt.Sprintf("There is no mapper available for items of type %q", itemType), parsedItem.PathProvider().GetWebRoute(item))
		return
	}

	panic("Unreachable")

}
