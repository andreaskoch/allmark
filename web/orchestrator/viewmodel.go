// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type ViewModelOrchestrator struct {
	*Orchestrator

	navigationOrchestrator NavigationOrchestrator
	tagOrchestrator        TagsOrchestrator
	fileOrchestrator       FileOrchestrator
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {

	// get the root item
	root := orchestrator.rootItem()
	if root == nil {
		orchestrator.logger.Warn("Cannot get viewmodel for route %q because no root item was found.", itemRoute)
		return viewModel, false
	}

	// get the requested item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		orchestrator.logger.Warn("Cannot get viewmodel for route %q because no item was found.", itemRoute)
		return viewModel, false
	}

	// convert content
	convertedContent, err := orchestrator.converter.Convert(orchestrator.relativePather(itemRoute), item)
	if err != nil {
		return viewModel, false
	}

	// create a view model
	viewModel = viewmodel.Model{
		Base:    getBaseModel(root, item, orchestrator.itemPather()),
		Content: convertedContent,
		Childs:  orchestrator.getChildModels(itemRoute, orchestrator.relativePather(itemRoute)),

		// navigation
		ToplevelNavigation:   orchestrator.navigationOrchestrator.GetToplevelNavigation(),
		BreadcrumbNavigation: orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(itemRoute),

		// tags
		Tags: orchestrator.tagOrchestrator.getItemTags(itemRoute),

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(itemRoute),
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

	return viewModel, true
}

func (orchestrator *ViewModelOrchestrator) getChildModels(itemRoute route.Route, pathProvider paths.Pather) []*viewmodel.Base {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	childModels := make([]*viewmodel.Base, 0)

	childItems := orchestrator.getChilds(itemRoute)
	for _, childItem := range childItems {
		baseModel := getBaseModel(rootItem, childItem, pathProvider)
		childModels = append(childModels, &baseModel)
	}

	return childModels
}

// sort the models by date and name
func sortModelsByDate(model1, model2 *viewmodel.Model) bool {

	return model1.CreationDate > model2.CreationDate

}
