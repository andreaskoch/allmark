// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"strings"
	"time"

	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/web/view/viewmodel"
)

func getBaseModel(root, item *model.Item, config config.Config) viewmodel.Base {

	baseModel := viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   item.Route().Value(),
		Level:   item.Route().Level(),
		BaseURL: GetBaseURL(item.Route()),
		Aliases: getAliasViewModels(item),

		PrintURL:    GetTypedItemURL(item.Route(), "print"),
		JSONURL:     GetTypedItemURL(item.Route(), "json"),
		MarkdownURL: GetTypedItemURL(item.Route(), "markdown"),

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
			baseModel.ParentRoute = parentRoute.Value()
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
func sortBaseModelsByDate(model1, model2 viewmodel.Base) bool {

	return model1.CreationDate > model2.CreationDate

}

// sort the models by date and name
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.After(model2.MetaData.CreationDate)

}

func pagedViewmodels(viewmodels []viewmodel.Model, pageSize, page int) (latest []viewmodel.Model, found bool) {

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(viewmodels) {
		return []viewmodel.Model{}, false
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

// getAliasViewModels returns a list of alias view-models for each alias of the specified item.
func getAliasViewModels(item *model.Item) []viewmodel.Alias {
	var viewModels []viewmodel.Alias
	for _, alias := range item.MetaData.Aliases {
		viewModels = append(viewModels, viewmodel.Alias{
			Name:        alias,
			Route:       "!" + alias, // Todo: Don't use a magic string for the alias prefix
			TargetRoute: item.Route().Value(),
		})
	}
	return viewModels
}
