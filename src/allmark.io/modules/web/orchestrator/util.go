// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"strings"

	"allmark.io/modules/common/config"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
)

func getBaseModel(root, item *model.Item, pathProvider paths.Pather) viewmodel.Base {

	baseModel := viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   pathProvider.Path(item.Route().Value()),
		Level:   item.Route().Level(),
		BaseUrl: GetBaseUrl(item.Route()),
		Alias:   item.MetaData.Alias,

		PrintUrl: GetTypedItemUrl(item.Route(), "print"),
		JsonUrl:  GetTypedItemUrl(item.Route(), "json"),

		PageTitle:   getPageTitleForItem(root, item),
		Title:       item.Title,
		Description: item.Description,

		LanguageTag:      getLanguageCode(item.MetaData.Language),
		CreationDate:     item.MetaData.CreationDate.Format("2006-01-02"),
		LastModifiedDate: item.MetaData.LastModifiedDate.Format("2006-01-02"),
	}

	if item.Route().Level() > 0 {
		if parentRoute, exists := item.Route().Parent(); exists {
			baseModel.ParentRoute = pathProvider.Path(parentRoute.Value())
		}
	}

	return baseModel

}

func getLanguageCode(languageHint string) string {
	if languageHint == "" {
		return config.DefaultLanguage
	}

	return languageHint
}

func getPageTitleForItem(rootItem, item *model.Item) string {
	if item.Route().Value() == rootItem.Route().Value() {
		return item.Title
	}

	return fmt.Sprintf("%s - %s", item.Title, rootItem.Title)
}

func GetBaseUrl(route route.Route) string {
	url := route.Value()
	if url != "" {
		return "/" + url + "/"
	}

	return "/"
}

func GetTypedItemUrl(route route.Route, urlType string) string {
	itemPath := GetBaseUrl(route)
	itemPath = strings.TrimSuffix(itemPath, "/")

	if len(itemPath) > 0 {
		return fmt.Sprintf("%s.%s", itemPath, urlType)
	}

	return urlType
}

// sort the models by date and name
func sortBaseModelsByDate(model1, model2 *viewmodel.Base) bool {

	return model1.CreationDate > model2.CreationDate

}

// sort the models by date and name
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.After(model2.MetaData.CreationDate)

}

func pagedViewmodels(viewmodels []*viewmodel.Model, pageSize, page int) (latest []*viewmodel.Model, found bool) {

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(viewmodels) {
		return []*viewmodel.Model{}, false
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(viewmodels) {
		endIndex = len(viewmodels)
	}

	return viewmodels[startIndex:endIndex], true
}

func pagedItems(models []*model.Item, pageSize, page int) (items []*model.Item, found bool) {

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(models) || page < 0 {
		return []*model.Item{}, false
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(models) {
		endIndex = len(models)
	}

	return models[startIndex:endIndex], true
}
