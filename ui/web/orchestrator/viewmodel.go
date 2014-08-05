// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewViewModelOrchestrator(itemIndex *index.Index, converter conversion.Converter, navigationOrchestrator *NavigationOrchestrator, tagOrchestrator *TagsOrchestrator) ViewModelOrchestrator {
	return ViewModelOrchestrator{
		itemIndex:              itemIndex,
		converter:              converter,
		navigationOrchestrator: navigationOrchestrator,
		tagOrchestrator:        tagOrchestrator,
		fileOrchestrator:       NewFileOrchestrator(),
	}
}

type ViewModelOrchestrator struct {
	itemIndex              *index.Index
	converter              conversion.Converter
	navigationOrchestrator *NavigationOrchestrator
	tagOrchestrator        *TagsOrchestrator
	fileOrchestrator       FileOrchestrator
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(pathProvider paths.Pather, item *model.Item) viewmodel.Model {

	// get the root item
	root := orchestrator.itemIndex.Root()
	if root == nil {
		return viewmodel.Model{}
	}

	// convert content
	convertedContent, err := orchestrator.converter.Convert(pathProvider, item, false)
	if err != nil {
		return viewmodel.Model{}
	}

	// create a view model
	viewModel := viewmodel.Model{
		Base:    getBaseModel(root, item, pathProvider),
		Content: convertedContent,
		Childs:  orchestrator.getChildModels(item.Route(), pathProvider),

		// navigation
		ToplevelNavigation:   orchestrator.navigationOrchestrator.GetToplevelNavigation(),
		BreadcrumbNavigation: orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(item),

		// tags
		Tags: orchestrator.tagOrchestrator.GetItemTags(item),

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(pathProvider, item),
	}

	// special viewmodel attributes
	isRepositoryItem := item.Type == model.TypeRepository
	if isRepositoryItem {

		// tag cloud
		repositoryIsNotEmpty := orchestrator.itemIndex.Size() > 5 // don't bother to create a tag cloud if there aren't enough documents
		if repositoryIsNotEmpty {

			tagCloud := orchestrator.tagOrchestrator.GetTagCloud()
			viewModel.TagCloud = tagCloud

		}

	}

	return viewModel
}

func (orchestrator *ViewModelOrchestrator) getChildModels(route *route.Route, pathProvider paths.Pather) []*viewmodel.Base {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	childModels := make([]*viewmodel.Base, 0)
	childItems := orchestrator.itemIndex.GetDirectChilds(route)
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
