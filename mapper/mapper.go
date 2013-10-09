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

func Map(item *repository.Item, itemResolver func(itemName string) *repository.Item, tagPath func(tag *repository.Tag) string, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) *view.Model {

	var model *view.Model

	// map the parsed item to the view model depending on the item type
	switch itemType := item.MetaData.ItemType; itemType {

	case types.PresentationItemType, types.RepositoryItemType, types.DocumentItemType, types.MessageItemType, types.LocationItemType:
		model = getModel(item, itemResolver, tagPath, relativePath, absolutePath, content)
		model.Childs = getSubModels(item, itemResolver, tagPath, relativePath, absolutePath, content)

	default:
		model = view.Error("Item type not recognized", fmt.Sprintf("There is no mapper available for items of type %q", itemType), relativePath(item), absolutePath(item))
	}

	// assign the model to the item
	item.Model = model

	return model
}

func getModel(item *repository.Item, itemResolver func(itemName string) *repository.Item, tagPath func(tag *repository.Tag) string, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) *view.Model {

	return &view.Model{
		Level:            item.Level,
		RelativeRoute:    relativePath(item),
		AbsoluteRoute:    absolutePath(item),
		Title:            item.Title,
		Description:      item.Description,
		Content:          content(item),
		LanguageTag:      getTwoLetterLanguageCode(item.MetaData.Language),
		CreationDate:     formatDate(item.MetaData.CreationDate),
		LastModifiedDate: formatDate(item.MetaData.LastModifiedDate),
		Type:             item.MetaData.ItemType,
		Tags:             getTags(item, tagPath),
		Locations:        getLocations(item.MetaData.Locations, itemResolver, tagPath, relativePath, absolutePath, content),
		GeoLocation:      getGeoLocation(item),
	}

}

func getSubModels(item *repository.Item, itemResolver func(itemName string) *repository.Item, tagPath func(tag *repository.Tag) string, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) []*view.Model {

	items := item.Childs
	models := make([]*view.Model, 0)

	for _, child := range items {
		models = append(models, Map(child, itemResolver, tagPath, relativePath, absolutePath, content))
	}

	return models
}

func getTags(item *repository.Item, tagPath func(tag *repository.Tag) string) []*view.Tag {
	tagModels := make([]*view.Tag, 0)

	for _, tag := range item.MetaData.Tags {
		tagModels = append(tagModels, &view.Tag{
			Name:          tag.Name(),
			AbsoluteRoute: tagPath(&tag),
		})
	}

	return tagModels
}
