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

func NewViewModelOrchestrator(itemIndex *index.ItemIndex, converter conversion.Converter, navigationOrchestrator *NavigationOrchestrator, tagOrchestrator *TagsOrchestrator) ViewModelOrchestrator {
	return ViewModelOrchestrator{
		itemIndex:              itemIndex,
		converter:              converter,
		navigationOrchestrator: navigationOrchestrator,
		tagOrchestrator:        tagOrchestrator,
	}
}

type ViewModelOrchestrator struct {
	itemIndex              *index.ItemIndex
	converter              conversion.Converter
	navigationOrchestrator *NavigationOrchestrator
	tagOrchestrator        *TagsOrchestrator
}

func (orchestrator *ViewModelOrchestrator) GetViewModel(pathProvider paths.Pather, item *model.Item) viewmodel.Model {

	// convert content
	convertedContent, err := orchestrator.converter.Convert(pathProvider, item)
	if err != nil {
		return viewmodel.Model{}
	}

	// create a view model
	viewModel := viewmodel.Model{
		Base:    getBaseModel(item, pathProvider),
		Content: convertedContent,
		Childs:  orchestrator.getChildModels(item.Route(), pathProvider),

		// navigation
		ToplevelNavigation:   orchestrator.navigationOrchestrator.GetToplevelNavigation(),
		BreadcrumbNavigation: orchestrator.navigationOrchestrator.GetBreadcrumbNavigation(item),

		// documents
		TopDocs: orchestrator.getTopDocuments(5, pathProvider, item.Route()),

		// tags
		Tags:     orchestrator.tagOrchestrator.GetItemTags(item),
		TagCloud: orchestrator.tagOrchestrator.GetTagCloud(),
	}

	return viewModel
}

func (orchestrator *ViewModelOrchestrator) getTopDocuments(numberOfTopDocuments int, pathProvider paths.Pather, route *route.Route) []*viewmodel.Model {

	baseRouteLevel := route.Level()
	nextRouteLevel := baseRouteLevel + 1
	childItems := orchestrator.itemIndex.GetAllChilds(route)

	// determine the candidates for the top-documents
	candidateModels := make([]*viewmodel.Model, 0)

	for len(candidateModels) == 0 && nextRouteLevel != baseRouteLevel+3 {

		for _, childItem := range childItems {

			if childItem.Route().Level() != nextRouteLevel {
				continue
			}

			// create viewmodel and append to list
			childModel := orchestrator.GetViewModel(pathProvider, childItem)
			candidateModels = append(candidateModels, &childModel)

		}

		nextRouteLevel++

	}

	// take the top models only
	if len(candidateModels) <= numberOfTopDocuments {
		return candidateModels
	}

	return candidateModels[:numberOfTopDocuments]
}

func (orchestrator *ViewModelOrchestrator) getChildModels(route *route.Route, pathProvider paths.Pather) []*viewmodel.Base {
	childModels := make([]*viewmodel.Base, 0)

	childItems := orchestrator.itemIndex.GetChilds(route)
	for _, childItem := range childItems {
		baseModel := getBaseModel(childItem, pathProvider)
		childModels = append(childModels, &baseModel)
	}

	return childModels
}
