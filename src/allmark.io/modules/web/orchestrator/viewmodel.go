// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
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

	// caches (do not initialize!)
	latestByRoute     map[string][]*viewmodel.Model
	viewmodelsByRoute map[string]*viewmodel.Model
}

func (orchestrator *ViewModelOrchestrator) GetFullViewModel(itemRoute route.Route) (viewmodel.Model, bool) {

	startTime := time.Now()

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return viewmodel.Model{}, false
	}

	// get the base view model
	viewModel := orchestrator.getViewModel(itemRoute)
	if viewModel == nil {
		return viewmodel.Model{}, false
	}

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
	orchestrator.logger.Statistics("Getting the full view model for route %q took %f seconds.", viewModel.Route, duration.Seconds())

	return *viewModel, true
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	vm := orchestrator.getViewModel(itemRoute)
	if vm == nil {
		return viewModel, false
	}

	return *vm, true
}

// GetViewModelByAlias returns the viewmodel by its alias.
func (orchestrator *ViewModelOrchestrator) GetViewModelByAlias(alias string) (viewModel viewmodel.Model, found bool) {

	item := orchestrator.getItemByAlias(alias)
	if item == nil {
		return viewmodel.Model{}, false
	}

	vm := orchestrator.getViewModel(item.Route())
	if vm == nil {
		return viewModel, false
	}

	return *vm, true
}

// GetLatest returns the latest items (sorted by creation date) for the given route.
func (orchestrator *ViewModelOrchestrator) GetLatest(itemRoute route.Route, pageSize, page int) (latest []*viewmodel.Model, found bool) {

	// return from cache if cache has been initialized
	if orchestrator.latestByRoute != nil {

		if models, exists := orchestrator.latestByRoute[itemRoute.Value()]; exists {
			return pagedViewmodels(models, pageSize, page)
		}

		return []*viewmodel.Model{}, false

	}

	// updateLatest updates the latest items for the given route.
	updateLatest := func(route route.Route) {
		startTime := time.Now()

		orchestrator.latestByRoute = make(map[string][]*viewmodel.Model)
		for _, childRoute := range orchestrator.repository.Routes() {
			latestItems := orchestrator.getLatestItems(childRoute)
			orchestrator.latestByRoute[childRoute.Value()] = orchestrator.getLastesViewModelsFromItemList(latestItems)
		}

		// log timing reports
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		orchestrator.logger.Statistics("Priming the latest cache took %f seconds.", duration.Seconds())
	}

	// asyncUpdateLatest executes updateLatest for the given route in a go routine.
	asyncUpdateLatest := func(route route.Route) {
		go updateLatest(route)
	}

	// initialize cache
	updateLatest(route.New())

	// register update callbacks
	orchestrator.registerUpdateCallback("update latest", UpdateTypeNew, asyncUpdateLatest)
	orchestrator.registerUpdateCallback("update latest", UpdateTypeModified, asyncUpdateLatest)
	orchestrator.registerUpdateCallback("update latest", UpdateTypeDeleted, asyncUpdateLatest)

	// return the result
	return orchestrator.GetLatest(itemRoute, pageSize, page)
}

// Converts a list of model.Item elements into a view models for the latest-items controller
func (orchestrator *ViewModelOrchestrator) getLastesViewModelsFromItemList(items []*model.Item) []*viewmodel.Model {

	// create viewmodels
	models := make([]*viewmodel.Model, 0, len(items))
	for _, item := range items {

		viewModel := orchestrator.getViewModel(item.Route())
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

func (orchestrator *ViewModelOrchestrator) getViewModel(itemRoute route.Route) *viewmodel.Model {

	if orchestrator.viewmodelsByRoute != nil {
		return orchestrator.viewmodelsByRoute[itemRoute.String()]
	}

	// updateViewModel stores the view model for the given route to the cache
	updateViewModel := func(route route.Route) {

		// convert content
		item := orchestrator.getItem(route)
		if item == nil {
			orchestrator.logger.Warn("Cannot update viewmodel cache. The item with the route %q was not found.", route.String())
			return
		}

		root := orchestrator.rootItem()
		convertedContent, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, orchestrator.relativePather(route), item)
		if err != nil {
			orchestrator.logger.Warn("Cannot convert content for route %q. Error: %s.", route, err.Error())
			convertedContent = "<!-- Conversion Error -->"
		}

		viewModel := &viewmodel.Model{
			Base:             getBaseModel(root, item, orchestrator.itemPather(), orchestrator.config),
			Content:          convertedContent,
			Publisher:        orchestrator.getPublisherInformation(),
			Author:           orchestrator.getAuthorInformation(item.MetaData.Author),
			Files:            orchestrator.fileOrchestrator.GetFiles(route),
			Images:           orchestrator.fileOrchestrator.GetImages(route),
			IsRepositoryItem: true,
		}

		// add rft url if rtf conversion is enabled
		if orchestrator.config.Conversion.Rtf.IsEnabled() {
			viewModel.RtfUrl = GetTypedItemUrl(route, "rtf")
		}

		orchestrator.viewmodelsByRoute[route.String()] = viewModel
	}

	// deleteViewModel stores the view model for the given route from the cache
	deleteViewModel := func(route route.Route) {
		delete(orchestrator.viewmodelsByRoute, route.String())
	}

	// build the cache
	orchestrator.viewmodelsByRoute = make(map[string]*viewmodel.Model)
	for _, item := range orchestrator.index().GetAllItems() {
		updateViewModel(item.Route())
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update viewmodel", UpdateTypeNew, updateViewModel)
	orchestrator.registerUpdateCallback("update viewmodel", UpdateTypeModified, updateViewModel)
	orchestrator.registerUpdateCallback("update viewmodel", UpdateTypeDeleted, deleteViewModel)

	return orchestrator.viewmodelsByRoute[itemRoute.String()]
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
		baseModel := getBaseModel(rootItem, childItem, pathProvider, orchestrator.config)
		childModels = append(childModels, &baseModel)
	}

	// sort the models
	viewmodel.SortBaseModelBy(sortBaseModelsByDate).Sort(childModels)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	orchestrator.logger.Statistics("Getting child models for route %q took %f seconds.", itemRoute, duration.Seconds())

	return childModels
}
