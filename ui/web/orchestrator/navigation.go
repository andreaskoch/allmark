// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewNavigationOrchestrator(itemIndex *index.Index, pathProvider paths.Pather) NavigationOrchestrator {
	return NavigationOrchestrator{
		itemIndex:    itemIndex,
		pathProvider: pathProvider,
	}
}

type NavigationOrchestrator struct {
	itemIndex    *index.Index
	pathProvider paths.Pather
}

func (orchestrator *NavigationOrchestrator) GetToplevelNavigation() *viewmodel.ToplevelNavigation {

	root := route.New()
	toplevelEntries := make([]*viewmodel.ToplevelEntry, 0)
	for _, child := range orchestrator.itemIndex.GetDirectChilds(root) {

		toplevelEntries = append(toplevelEntries, &viewmodel.ToplevelEntry{
			Title: child.Title,
			Path:  orchestrator.pathProvider.Path(child.Route().Value()),
		})

	}

	return &viewmodel.ToplevelNavigation{
		Entries: toplevelEntries,
	}
}

func (orchestrator *NavigationOrchestrator) GetBreadcrumbNavigation(item *model.Item) *viewmodel.BreadcrumbNavigation {

	// create a new bread crumb navigation
	navigation := &viewmodel.BreadcrumbNavigation{
		Entries: make([]*viewmodel.Breadcrumb, 0),
	}

	// abort if item or model is nil
	if item == nil {
		return navigation
	}

	// recurse if there is a parent
	if parent := orchestrator.itemIndex.GetParent(item.Route()); parent != nil {
		navigation.Entries = append(navigation.Entries, orchestrator.GetBreadcrumbNavigation(parent).Entries...)
	}

	// append a new navigation entry and return it
	navigation.Entries = append(navigation.Entries, &viewmodel.Breadcrumb{
		Title: item.Title,
		Level: item.Route().Level(),
		Path:  orchestrator.pathProvider.Path(item.Route().Value()),
	})

	return navigation
}
