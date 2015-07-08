// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"strings"
	"time"

	"allmark.io/modules/common/config"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
)

func getBaseModel(root, item *model.Item, pathProvider paths.Pather, config config.Config) viewmodel.Base {

	baseModel := viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   pathProvider.Path(item.Route().Value()),
		Level:   item.Route().Level(),
		BaseURL: GetBaseURL(item.Route()),
		Alias:   item.MetaData.Alias,

		PrintURL: GetTypedItemURL(item.Route(), "print"),
		JsonURL:  GetTypedItemURL(item.Route(), "json"),

		PageTitle:   getPageTitleForItem(root, item),
		Title:       item.Title,
		Description: item.Description,

		LanguageTag:      getLanguageCode(item.MetaData.Language),
		CreationDate:     getFormattedDate(item.MetaData.CreationDate),
		LastModifiedDate: getFormattedDate(item.MetaData.LastModifiedDate),

		LiveReloadEnabled: config.LiveReload.Enabled,
	}

	if item.Route().Level() > 0 {
		if parentRoute, exists := item.Route().Parent(); exists {
			baseModel.ParentRoute = pathProvider.Path(parentRoute.Value())
		}
	}

	return baseModel

}

// Get the formatted date if the supplied date is initialized; otherwise return an empty string.
func getFormattedDate(date time.Time) string {

	if date.IsZero() {
		return ""
	}

	return date.Format("2006-01-02")
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

func GetBaseURL(route route.Route) string {
	url := route.Value()
	if url != "" {
		return "/" + url + "/"
	}

	return "/"
}

func GetTypedItemURL(route route.Route, urlType string) string {
	itemPath := GetBaseURL(route)
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
