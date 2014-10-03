// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type NavigationOrchestrator struct {
	*Orchestrator
}

func (orchestrator *NavigationOrchestrator) GetToplevelNavigation() *viewmodel.ToplevelNavigation {

	root := route.New()
	toplevelEntries := make([]*viewmodel.ToplevelEntry, 0)

	for _, child := range orchestrator.getChilds(root) {

		toplevelEntries = append(toplevelEntries, &viewmodel.ToplevelEntry{
			Title: child.Title,
			Path:  orchestrator.itemPather().Path(child.Route().Value()),
		})

	}

	return &viewmodel.ToplevelNavigation{
		Entries: toplevelEntries,
	}
}

func (orchestrator *NavigationOrchestrator) GetBreadcrumbNavigation(route route.Route) *viewmodel.BreadcrumbNavigation {

	// create a new bread crumb navigation
	navigation := &viewmodel.BreadcrumbNavigation{
		Entries: make([]*viewmodel.Breadcrumb, 0),
	}

	// get the item for the supplied route
	item := orchestrator.getItem(route)
	if item == nil {
		orchestrator.logger.Debug("Returning an empty navigation model because there is no item for route %q.", route)
		return navigation
	}

	// recurse if there is a parent
	if parent := orchestrator.getParent(item.Route()); parent != nil {
		navigation.Entries = append(navigation.Entries, orchestrator.GetBreadcrumbNavigation(parent.Route()).Entries...)
	}

	// append a new navigation entry and return it
	navigation.Entries = append(navigation.Entries, &viewmodel.Breadcrumb{
		Title: item.Title,
		Level: item.Route().Level(),
		Path:  orchestrator.itemPather().Path(item.Route().Value()),
	})

	return navigation
}

func (orchestrator *NavigationOrchestrator) GetItemNavigation(route route.Route) *viewmodel.ItemNavigation {

	// create a new item navigation
	navigation := &viewmodel.ItemNavigation{}

	// get the item for the supplied route
	item := orchestrator.getItem(route)
	if item == nil {
		return navigation
	}

	// get the parent
	if route.Level() > 0 {
		if parent := orchestrator.getParent(item.Route()); parent != nil {
			navigation.Parent = &viewmodel.NavEntry{
				Title:       parent.Title,
				Description: parent.Description,
				Path:        orchestrator.itemPather().Path(parent.Route().Value()),
			}
		}
	}

	return navigation
}
