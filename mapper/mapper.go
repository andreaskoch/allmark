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

func Map(item *repository.Item) *view.Model {

	var model *view.Model

	// map the parsed item to the view model depending on the item type
	switch itemType := item.MetaData.ItemType; itemType {

	case types.PresentationItemType, types.RepositoryItemType, types.DocumentItemType, types.MessageItemType:
		model = getModel(item)
		model.Childs = getSubModels(item)

	default:
		model = view.Error("Item type not recognized", fmt.Sprintf("There is no mapper available for items of type %q", itemType), item.RelativePath, item.AbsolutePath)
	}

	// assign the model to the item
	item.Model = model

	return model
}

func getModel(item *repository.Item) *view.Model {

	return &view.Model{
		Level:         item.Level,
		RelativeRoute: item.RelativePath,
		AbsoluteRoute: item.AbsolutePath,
		Title:         item.Title,
		Description:   item.Description,
		LanguageTag:   getTwoLetterLanguageCode(item.MetaData.Language),
		Date:          formatDate(item.MetaData.Date),
		Type:          item.MetaData.ItemType,
	}

}

func getSubModels(item *repository.Item) []*view.Model {

	items := item.Childs
	models := make([]*view.Model, 0)

	for _, child := range items {
		models = append(models, Map(child))
	}

	return models
}
