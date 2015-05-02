// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"time"

	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	converter "allmark.io/modules/services/converter/markdowntohtml"
	"allmark.io/modules/web/view/viewmodel"
)

type ViewModelOrchestrator struct {
	*Orchestrator

	navigationOrchestrator *NavigationOrchestrator
	tagOrchestrator        *TagsOrchestrator
	fileOrchestrator       *FileOrchestrator

	latestByRoute     map[string][]*viewmodel.Model
	viewmodelsByRoute map[string]*viewmodel.Model
}

func (orchestrator *ViewModelOrchestrator) blockingCacheWarmup() {
	orchestrator.getViewModel(orchestrator.rootItem(), false)
	orchestrator.GetLatest(route.New(), 5, 1)
}

func (orchestrator *ViewModelOrchestrator) GetFullViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	startTime := time.Now()

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	// get the base view model
	viewModel = *orchestrator.getViewModel(item, false)

	// navigation
	viewModel.ToplevelNavigation = orchestrator.navigationOrchestrator.GetToplevelNavigation()
	viewModel.BreadcrumbNavigation = orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(itemRoute)
	viewModel.ItemNavigation = orchestrator.navigationOrchestrator.GetItemNavigation(itemRoute)

	// childs
	viewModel.Childs = orchestrator.getChildModels(itemRoute)

	// tags
	viewModel.Tags = orchestrator.tagOrchestrator.getItemTags(itemRoute)

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

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	orchestrator.logger.Statistics("Getting the full view model %s took %f seconds.", viewModel.Route, duration.Seconds())

	return viewModel, true
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewModel, false
	}

	return *orchestrator.getViewModel(item, false), true
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
			return pagedViewmodels(models, pageSize, page)
		}

		return []*viewmodel.Model{}, false

	}

	orchestrator.setCache(cacheType, func() {

		startTime := time.Now()

		latestModelsByRoute := make(map[string][]*viewmodel.Model)

		for _, childRoute := range orchestrator.repository.Routes() {

			// get the latest items
			latestItems := orchestrator.getLatestItems(childRoute)

			// store the results
			latestModelsByRoute[childRoute.Value()] = orchestrator.getLastesViewModelsFromItemList(latestItems)

		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		orchestrator.logger.Statistics("Priming the latest cache took %f seconds.", duration.Seconds())

		// save the result
		orchestrator.latestByRoute = latestModelsByRoute

	})

	// return a result
	if models, exists := orchestrator.latestByRoute[itemRoute.Value()]; exists {
		return pagedViewmodels(models, pageSize, page)
	}

	return []*viewmodel.Model{}, false
}

// Converts a list of model.Item elements into a view models for the latest-items controller
func (orchestrator *ViewModelOrchestrator) getLastesViewModelsFromItemList(items []*model.Item) []*viewmodel.Model {

	// create viewmodels
	models := make([]*viewmodel.Model, 0, len(items))
	for _, item := range items {

		viewModel := orchestrator.getViewModel(item, false)
		if viewModel == nil {
			orchestrator.logger.Error("No view model found for item %q.", item)
			continue
		}

		// prepare lazy loading
		contentWithLazyLoadingEnabled := converter.LazyLoad(viewModel.Content)

		// create a copy (make sure we don't modify the content of the original view model)
		viewmodelCopy := *viewModel
		viewmodelCopy.Content = contentWithLazyLoadingEnabled

		models = append(models, &viewmodelCopy)
	}

	return models
}

func (orchestrator *ViewModelOrchestrator) getViewModel(item *model.Item, skipCache bool) *viewmodel.Model {

	// get the root item
	root := orchestrator.rootItem()
	if root == nil {
		panic(fmt.Sprintf("Cannot get viewmodel for route %q because no root item was found.", item))
	}

	itemRoute := item.Route()
	cacheType := "viewmodels"

	if skipCache {
		// prime the cache synchronously
		orchestrator.primeCache(cacheType)
	}

	// load from cache
	if orchestrator.viewmodelsByRoute != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		// return the result
		return orchestrator.viewmodelsByRoute[itemRoute.Value()]
	}

	orchestrator.setCache(cacheType, func() {

		startTime := time.Now()

		viewmodelsByRoute := make(map[string]*viewmodel.Model)

		for _, child := range orchestrator.index().GetAllItems() {

			childRoute := child.Route()

			// convert content
			convertedContent, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, orchestrator.relativePather(childRoute), child)
			if err != nil {
				orchestrator.logger.Warn("Cannot convert content for item %q. Error: %s.", child.String(), err.Error())
				convertedContent = "<!-- Conversion Error -->"
			}

			// create a view model
			viewModel := &viewmodel.Model{
				Base:             getBaseModel(root, child, orchestrator.itemPather()),
				Content:          convertedContent,
				Publisher:        orchestrator.getPublisherInformation(),
				Author:           orchestrator.getAuthorInformation(child.MetaData.Author),
				Files:            orchestrator.fileOrchestrator.GetFiles(childRoute),
				Images:           orchestrator.fileOrchestrator.GetImages(childRoute),
				IsRepositoryItem: true,
			}

			// add rft url if rtf conversion is enabled
			if orchestrator.config.Conversion.Rtf.Enabled {
				viewModel.RtfUrl = GetTypedItemUrl(item.Route(), "rtf")
			}

			// store the view model
			viewmodelsByRoute[childRoute.Value()] = viewModel

		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		orchestrator.logger.Statistics("Priming the viewModel cache took %f seconds.", duration.Seconds())

		orchestrator.viewmodelsByRoute = viewmodelsByRoute

	})

	return orchestrator.viewmodelsByRoute[itemRoute.Value()]
}

func (orchestrator *ViewModelOrchestrator) getChildModels(itemRoute route.Route) []*viewmodel.Base {

	startTime := time.Now()

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

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	orchestrator.logger.Statistics("Getting child models for route %s took %f seconds.", itemRoute, duration.Seconds())

	return childModels
}
