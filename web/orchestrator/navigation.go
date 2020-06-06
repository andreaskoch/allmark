// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/web/view/viewmodel"
)

type NavigationOrchestrator struct {
	*Orchestrator

	toplevelNavigation *viewmodel.ToplevelNavigation
}

func (orchestrator *NavigationOrchestrator) GetToplevelNavigation() viewmodel.ToplevelNavigation {

	if orchestrator.toplevelNavigation != nil {
		return *orchestrator.toplevelNavigation
	}

	// updateToplevelNavigation creates a new toplevel navigation and stores it in the cache
	updateToplevelNavigation := func(r route.Route) {
		root := route.New()
		toplevelEntries := make([]viewmodel.ToplevelEntry, 0)

		for _, child := range orchestrator.getChildren(root) {

			toplevelEntries = append(toplevelEntries, viewmodel.ToplevelEntry{
				Title: child.Title,
				Path:  orchestrator.itemPather().Path(child.Route().Value()),
			})

		}

		orchestrator.toplevelNavigation = &viewmodel.ToplevelNavigation{
			Entries: toplevelEntries,
		}
	}

	// write the cache
	updateToplevelNavigation(route.New())

	// register update callbacks
	orchestrator.registerUpdateCallback("update toplevel navigation", UpdateTypeNew, updateToplevelNavigation)
	orchestrator.registerUpdateCallback("update toplevel navigation", UpdateTypeModified, updateToplevelNavigation)
	orchestrator.registerUpdateCallback("update toplevel navigation", UpdateTypeDeleted, updateToplevelNavigation)

	return orchestrator.GetToplevelNavigation()
}

func (orchestrator *NavigationOrchestrator) GetBreadcrumbNavigation(route route.Route) viewmodel.BreadcrumbNavigation {

	// create a new bread crumb navigation
	navigation := viewmodel.BreadcrumbNavigation{
		Entries: make([]viewmodel.Breadcrumb, 0),
	}

	// get the item for the supplied route
	item := orchestrator.getItem(route)
	if item == nil {
		orchestrator.logger.Error("Returning an empty navigation model because there is no item for route %q.", route)
		return navigation
	}

	// recurse if there is a parent
	if parent := orchestrator.getParent(item.Route()); parent != nil {
		navigation.Entries = append(navigation.Entries, orchestrator.GetBreadcrumbNavigation(parent.Route()).Entries...)
	}

	// append a new navigation entry and return it
	unmarkedEntries := append(navigation.Entries, viewmodel.Breadcrumb{
		Title: item.Title,
		Level: item.Route().Level(),
		Path:  orchestrator.itemPather().Path(item.Route().Value()),
	})

	// mark the entries
	markdedEntries := make([]viewmodel.Breadcrumb, 0)
	for index, entry := range unmarkedEntries {

		if index < (len(unmarkedEntries) - 1) {
			entry.IsLast = false
		} else {
			entry.IsLast = true
		}

		markdedEntries = append(markdedEntries, entry)
	}

	navigation.Entries = markdedEntries

	return navigation
}

func (orchestrator *NavigationOrchestrator) GetItemNavigation(route route.Route) viewmodel.ItemNavigation {

	// create a new item navigation
	navigation := viewmodel.ItemNavigation{}

	// get the item for the supplied route
	item := orchestrator.getItem(route)
	if item == nil {
		return navigation
	}

	// get the parent
	if parent := orchestrator.getParent(item.Route()); parent != nil {
		navigation.Parent = viewmodel.NavEntry{
			Title:       parent.Title,
			Description: parent.Description,
			Path:        orchestrator.itemPather().Path(parent.Route().Value()),
		}
	}

	// previous
	if previous := orchestrator.getPrevious(item.Route()); previous != nil {
		navigation.Previous = viewmodel.NavEntry{
			Title:       previous.Title,
			Description: previous.Description,
			Path:        orchestrator.itemPather().Path(previous.Route().Value()),
		}
	}

	// next
	if next := orchestrator.getNext(item.Route()); next != nil {
		navigation.Next = viewmodel.NavEntry{
			Title:       next.Title,
			Description: next.Description,
			Path:        orchestrator.itemPather().Path(next.Route().Value()),
		}
	}

	return navigation
}
