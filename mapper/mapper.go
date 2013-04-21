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

type Mapper func(item *repository.Item) view.Model

func New(itemType string, pathProvider *path.Provider, targetFormat string) Mapper {

	switch itemType {
	case repository.DocumentItemType:
		return createDocumentMapperFunc(pathProvider, targetFormat)

	case repository.MessageItemType:
		return createMessageMapperFunc(pathProvider, targetFormat)

	case repository.RepositoryItemType, repository.CollectionItemType:
		return createCollectionMapperFunc(pathProvider, targetFormat)
	}

	return func(item *repository.Item) view.Model {
		return view.Error(fmt.Sprintf("There is no mapper available for items of type %q", itemType), pathProvider.GetWebRoute(item))
	}
}
