// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"sort"
	"strings"
	"time"
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
		RtfUrl:   GetTypedItemUrl(item.Route(), "rtf"),

		PageTitle:   getPageTitle(root, item),
		Title:       item.Title,
		Description: item.Description,

		LanguageTag:      item.MetaData.Language,
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

func getPageTitle(rootItem, item *model.Item) string {
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
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.After(model2.MetaData.CreationDate)

}

// sort the models by date and name
func sortBaseModelsByDate(model1, model2 *viewmodel.Base) bool {

	return model1.CreationDate > model2.CreationDate

}

func sortRoutesAndDatesDescending(itemRoute1, itemRoute2 routeAndDate) bool {
	return itemRoute1.date.After(itemRoute2.date)
}

type routeAndDate struct {
	route route.Route
	date  time.Time
}

type SortItemRoutesAndDatesBy func(itemRoute1, itemRoute2 routeAndDate) bool

func (by SortItemRoutesAndDatesBy) Sort(routesAndDates []routeAndDate) {
	sorter := &routeAndDateSorter{
		routesAndDates: routesAndDates,
		by:             by,
	}

	sort.Sort(sorter)
}

type routeAndDateSorter struct {
	routesAndDates []routeAndDate
	by             SortItemRoutesAndDatesBy
}

func (sorter *routeAndDateSorter) Len() int {
	return len(sorter.routesAndDates)
}

func (sorter *routeAndDateSorter) Swap(i, j int) {
	sorter.routesAndDates[i], sorter.routesAndDates[j] = sorter.routesAndDates[j], sorter.routesAndDates[i]
}

func (sorter *routeAndDateSorter) Less(i, j int) bool {
	return sorter.by(sorter.routesAndDates[i], sorter.routesAndDates[j])
}
