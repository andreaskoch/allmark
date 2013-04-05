// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func GetMapper(pathProvider *path.Provider, item *repository.Item) (func(item *repository.Item) view.Model, error) {

	switch item.Type {
	case repository.DocumentItemType:
		return createDocumentMapperFunc(pathProvider), nil

	case repository.MessageItemType:
		return createMessageMapperFunc(pathProvider), nil

	case repository.RepositoryItemType, repository.CollectionItemType:
		return createCollectionMapperFunc(pathProvider), nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))
}
