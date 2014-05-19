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
	convertedContent, err := orchestrator.converter.Convert(pathProvider, item)
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

	// attach tag cloud and top documents if the item is a repository type
	if item.Type == model.TypeRepository {
		viewModel.TopDocs = orchestrator.getTopDocuments(5, pathProvider, item.Route())
		viewModel.TagCloud = orchestrator.tagOrchestrator.GetTagCloud()
	}

	return viewModel
}

func (orchestrator *ViewModelOrchestrator) getTopDocuments(numberOfTopDocuments int, pathProvider paths.Pather, route *route.Route) []*viewmodel.Model {

	baseRouteLevel := route.Level()

	// determine the candidates for the top-documents
	candidateModels := make([]*viewmodel.Model, 0)

	// include only the next or over-next level childs
	nextLevelChildExpression := func(child *model.Item) bool {
		childLevel := child.Route().Level()

		isNextLevel := childLevel == baseRouteLevel+1
		isOverNextLevel := childLevel == baseRouteLevel+2

		return isNextLevel || isOverNextLevel
	}

	childs := orchestrator.itemIndex.GetAllChilds(route, nextLevelChildExpression)
	for _, child := range childs {

		// filter out virtual items
		if child.IsVirtual() {
			continue
		}

		// create viewmodel and append to list
		childModel := orchestrator.GetViewModel(pathProvider, child)
		candidateModels = append(candidateModels, &childModel)

	}

	// sort the candidate models
	viewmodel.SortModelBy(sortModelsByDate).Sort(candidateModels)

	// take the top models only
	if len(candidateModels) <= numberOfTopDocuments {
		return candidateModels
	}

	return candidateModels[:numberOfTopDocuments]
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
