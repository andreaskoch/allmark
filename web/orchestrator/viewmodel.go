// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	converter "github.com/andreaskoch/allmark2/services/converter/markdowntohtml"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type ViewModelOrchestrator struct {
	*Orchestrator

	navigationOrchestrator NavigationOrchestrator
	tagOrchestrator        TagsOrchestrator
	fileOrchestrator       FileOrchestrator
	locationOrchestrator   LocationOrchestrator
}

func (orchestrator *ViewModelOrchestrator) GetFullViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	// get the base view model
	viewModel = orchestrator.getViewModel(item)

	// navigation
	viewModel.ToplevelNavigation = orchestrator.navigationOrchestrator.GetToplevelNavigation()
	viewModel.BreadcrumbNavigation = orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(itemRoute)
	viewModel.ItemNavigation = orchestrator.navigationOrchestrator.GetItemNavigation(itemRoute)

	// childs
	viewModel.Childs = orchestrator.getChildModels(itemRoute)

	// tags
	viewModel.Tags = orchestrator.tagOrchestrator.getItemTags(itemRoute)

	// Locations
	viewModel.Locations = orchestrator.locationOrchestrator.GetLocations(item.MetaData.Locations, func(i *model.Item) viewmodel.Model {
		return orchestrator.getViewModel(i)
	})

	// Geo Coordinates
	viewModel.GeoLocation = getGeoLocation(item)

	// Analytics Settings
	viewModel.Analytics = orchestrator.getAnalyticsSettings()

	// Hash / ETag
	viewModel.Hash = item.Hash

	// special viewmodel attributes
	isRepositoryItem := item.Type == model.TypeRepository
	if isRepositoryItem {

		// tag cloud
		repositoryIsNotEmpty := orchestrator.index().Size() >= 5 // don't bother to create a tag cloud if there aren't enough documents
		if repositoryIsNotEmpty {

			tagCloud := orchestrator.tagOrchestrator.GetTagCloud()
			viewModel.TagCloud = tagCloud

		}

	}

	return viewModel, true
}

func (orchestrator *ViewModelOrchestrator) GetLatest(itemRoute route.Route, pageSize, page int) (models []*viewmodel.Model, found bool) {

	// get the latest routes
	latestRoutes, found := orchestrator.getLatestRoutesByPage(itemRoute, pageSize, page)
	if !found {
		return models, false
	}

	// create viewmodels
	models = make([]*viewmodel.Model, 0, len(latestRoutes))
	for _, route := range latestRoutes {

		viewModel, found := orchestrator.GetViewModel(route)
		if !found {
			orchestrator.logger.Warn("Viewmode %q was not found.", route)
			continue
		}

		// prepare lazy loading
		viewModel.Content = converter.LazyLoad(viewModel.Content)

		models = append(models, &viewModel)
	}

	return models, true
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	return orchestrator.getViewModel(item), true

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

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(itemRoute),
	}

	return viewModel
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
