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

	latestByRoute map[string][]*viewmodel.Model
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

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	return orchestrator.getViewModel(item), true

}

func (orchestrator *ViewModelOrchestrator) GetLatest(itemRoute route.Route, pageSize, page int) (latest []*viewmodel.Model, found bool) {

	cacheType := "latest"

	// load from cache
	if orchestrator.latestByRoute != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		// return the result
		if models, exists := orchestrator.latestByRoute[itemRoute.Value()]; exists {
			return pagedViewmodels(models, pageSize, page), true
		}

		return []*viewmodel.Model{}, false

	}

	orchestrator.setCache(cacheType, func() {

		latestModelsByRoute := make(map[string][]*viewmodel.Model)

		for _, childRoute := range orchestrator.repository.Routes() {

			// get the latest items
			latestItems := orchestrator.getLatestItems(childRoute)

			// create viewmodels
			models := make([]*viewmodel.Model, 0, len(latestItems))
			for _, item := range latestItems {

				viewModel := orchestrator.getViewModel(item)

				// prepare lazy loading
				viewModel.Content = converter.LazyLoad(viewModel.Content)

				models = append(models, &viewModel)
			}

			// store the results
			latestModelsByRoute[childRoute.Value()] = models

		}

		// save the result
		orchestrator.latestByRoute = latestModelsByRoute

	})

	// return a result
	if models, exists := orchestrator.latestByRoute[itemRoute.Value()]; exists {
		return pagedViewmodels(models, pageSize, page), true
	}

	return []*viewmodel.Model{}, false
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

		IsRepositoryItem: true,
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
