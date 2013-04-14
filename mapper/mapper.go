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

func GetMapper(pathProvider *path.Provider, converterFactory func(item *repository.Item) func() string, itemType string) (func(item *repository.Item) view.Model, error) {

	switch itemType {
	case repository.DocumentItemType:
		return createDocumentMapperFunc(pathProvider, converterFactory), nil

	case repository.MessageItemType:
		return createMessageMapperFunc(pathProvider, converterFactory), nil

	case repository.RepositoryItemType, repository.CollectionItemType:
		return createCollectionMapperFunc(pathProvider, converterFactory), nil
	}

	return nil, fmt.Errorf("There is no mapper available for items of type %q", itemType)
}
