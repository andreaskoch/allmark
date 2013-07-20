// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/view"
)

func MapItem(item *repository.Item) *view.Model {

	var model *view.Model

	// map the parsed item to the view model depending on the item type
	switch itemType := item.MetaData.ItemType; itemType {
	case types.PresentationItemType:
		model = createPresentationMapperFunc(item)

	case types.RepositoryItemType, types.DocumentItemType, types.MessageItemType:
		model = createDocumentMapperFunc(item)
		model.Childs = getSubModels(item)

	default:
		model = view.Error("Item type not recognized", fmt.Sprintf("There is no mapper available for items of type %q", itemType), item.RelativePath, item.AbsolutePath)
	}

	// assign the model to the item
	item.Model = model

	return model
}

func getSubModels(item *repository.Item) []*view.Model {

	items := item.Childs
	models := make([]*view.Model, 0)

	for _, child := range items {
		models = append(models, MapItem(child))
	}

	return models
}
