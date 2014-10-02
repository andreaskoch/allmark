// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"sort"
	"time"
)

type ViewModelOrchestrator struct {
	*Orchestrator

	navigationOrchestrator NavigationOrchestrator
	tagOrchestrator        TagsOrchestrator
	fileOrchestrator       FileOrchestrator
	locationOrchestrator   LocationOrchestrator

	leafesByRoute map[string][]route.Route
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	return orchestrator.getViewModel(item), true

}

func (orchestrator *ViewModelOrchestrator) GetLatest(itemRoute route.Route, pageSize, page int) (models []*viewmodel.Model, found bool) {

	leafes := orchestrator.getAllLeafes(itemRoute)

	// collect the creation dates for all leafes
	routesAndDates := make([]routeAndDate, len(leafes))
	for _, leaf := range leafes {
		creationDate, found := orchestrator.getCreationDate(leaf)
		if !found {
			// todo: log info
			continue
		}

		routesAndDates = append(routesAndDates, routeAndDate{leaf, creationDate})
	}

	// sort the leafes by date
	SortItemRoutesAndDatesBy(sortRoutesAndDatesDescending).Sort(routesAndDates)

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(routesAndDates) {
		return models, false
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(routesAndDates) {
		endIndex = len(routesAndDates)
	}

	selectedRoutesAndDates := routesAndDates[startIndex:endIndex]
	models = make([]*viewmodel.Model, len(selectedRoutesAndDates))
	for _, itemRoute := range selectedRoutesAndDates {

		viewModel, found := orchestrator.GetViewModel(itemRoute.route)
		if !found {
			// todo: log error
			continue
		}

		models = append(models, &viewModel)
	}

	return models, true
}

func (orchestrator *ViewModelOrchestrator) getViewModel(item *model.Item) viewmodel.Model {

	itemRoute := item.Route()

	// get the root item
	root := orchestrator.rootItem()
	if root == nil {
		panic(fmt.Sprintf("Cannot get viewmodel for route %q because no root item was found.", itemRoute))
	}

	// convert content
	convertedContent, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, orchestrator.relativePather(itemRoute), item)
	if err != nil {
		orchestrator.logger.Warn("Cannot convert content for item %q. Error: %s.", item.String(), err.Error())
		convertedContent = "<!-- Conversion Error -->"
	}

	// create a view model
	viewModel := viewmodel.Model{
		Base:    getBaseModel(root, item, orchestrator.itemPather()),
		Content: convertedContent,
		Childs:  orchestrator.getChildModels(itemRoute),

		// navigation
		ToplevelNavigation:   orchestrator.navigationOrchestrator.GetToplevelNavigation(),
		BreadcrumbNavigation: orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(itemRoute),

		// tags
		Tags: orchestrator.tagOrchestrator.getItemTags(itemRoute),

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(itemRoute),

		// Locations
		Locations: orchestrator.locationOrchestrator.GetLocations(item.MetaData.Locations, func(i *model.Item) viewmodel.Model {
			return orchestrator.getViewModel(i)
		}),

		// Geo Coordinates
		GeoLocation: getGeoLocation(item),
	}

	// special viewmodel attributes
	isRepositoryItem := item.Type == model.TypeRepository
	if isRepositoryItem {

		// tag cloud
		repositoryIsNotEmpty := orchestrator.repository.Size() > 5 // don't bother to create a tag cloud if there aren't enough documents
		if repositoryIsNotEmpty {

			tagCloud := orchestrator.tagOrchestrator.GetTagCloud()
			viewModel.TagCloud = tagCloud

		}

	}

	return viewModel
}

func (orchestrator *ViewModelOrchestrator) getAllLeafes(parentRoute route.Route) []route.Route {

	// cache lookup
	key := parentRoute.Value()
	if leafes, isset := orchestrator.leafesByRoute[key]; isset {
		return leafes
	}

	childRoutes := make([]route.Route, 0)

	childItems := orchestrator.getChilds(parentRoute)
	if hasNoMoreChilds := len(childItems) == 0; hasNoMoreChilds {
		return []route.Route{parentRoute}
	}

	// recurse
	for _, childItem := range childItems {
		childRoutes = append(childRoutes, orchestrator.getAllLeafes(childItem.Route())...)
	}

	// store the value
	orchestrator.leafesByRoute[key] = childRoutes

	return childRoutes

}

func (orchestrator *ViewModelOrchestrator) getChildModels(itemRoute route.Route) []*viewmodel.Base {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	pathProvider := orchestrator.relativePather(itemRoute)

	childModels := make([]*viewmodel.Base, 0)

	childItems := orchestrator.getChilds(itemRoute)
	for _, childItem := range childItems {
		baseModel := getBaseModel(rootItem, childItem, pathProvider)
		childModels = append(childModels, &baseModel)
	}

	// sort the models
	viewmodel.SortBaseModelBy(sortBaseModelsByDate).Sort(childModels)

	return childModels
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
