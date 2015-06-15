// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
)

type NavigationOrchestrator struct {
	*Orchestrator

	breadcrumbNavigationByRoute map[string]*viewmodel.BreadcrumbNavigation
	toplevelNavigation          *viewmodel.ToplevelNavigation
}

func (orchestrator *NavigationOrchestrator) GetToplevelNavigation() *viewmodel.ToplevelNavigation {

	if orchestrator.toplevelNavigation != nil {
		return orchestrator.toplevelNavigation
	}

	// updateToplevelNavigation creates a new toplevel navigation and stores it in the cache
	updateToplevelNavigation := func(r route.Route) {
		root := route.New()
		toplevelEntries := make([]*viewmodel.ToplevelEntry, 0)

		for _, child := range orchestrator.getChilds(root) {

			toplevelEntries = append(toplevelEntries, &viewmodel.ToplevelEntry{
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

func (orchestrator *NavigationOrchestrator) GetBreadcrumbNavigation(itemRoute route.Route) *viewmodel.BreadcrumbNavigation {

	// return from cache if possible
	if orchestrator.breadcrumbNavigationByRoute != nil {
		return orchestrator.breadcrumbNavigationByRoute[itemRoute.String()]
	}

	// initialize the cache
	orchestrator.breadcrumbNavigationByRoute = make(map[string]*viewmodel.BreadcrumbNavigation)

	// updateBreadcrumbNavigation writes a breadcrumb navigation model for the given route into the cache.
	updateBreadcrumbNavigation := func(route route.Route) {

		// create a new bread crumb navigation
		navigation := &viewmodel.BreadcrumbNavigation{
			Entries: make([]*viewmodel.Breadcrumb, 0),
		}

		// get the item for the supplied route
		item := orchestrator.getItem(route)
		if item == nil {
			orchestrator.logger.Error("Returning an empty navigation model because there is no item for route %q.", route)
			return
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

		// mark the entries
		for index, entry := range navigation.Entries {
			if index < (len(navigation.Entries) - 1) {
				entry.IsLast = false
			} else {
				entry.IsLast = true
			}
		}

		orchestrator.breadcrumbNavigationByRoute[item.Route().String()] = navigation
	}

	// deleteBreadcrumbNavigation removes the cache entry for the given route.
	deleteBreadcrumbNavigation := func(route route.Route) {
		delete(orchestrator.breadcrumbNavigationByRoute, route.String())
	}

	// write cache for all routes
	for _, childRoute := range orchestrator.repository.Routes() {
		updateBreadcrumbNavigation(childRoute)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update breadcrumb navigation", UpdateTypeNew, updateBreadcrumbNavigation)
	orchestrator.registerUpdateCallback("update breadcrumb navigation", UpdateTypeModified, updateBreadcrumbNavigation)
	orchestrator.registerUpdateCallback("update breadcrumb navigation", UpdateTypeDeleted, deleteBreadcrumbNavigation)

	return orchestrator.GetBreadcrumbNavigation(itemRoute)
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
	if parent := orchestrator.getParent(item.Route()); parent != nil {
		navigation.Parent = &viewmodel.NavEntry{
			Title:       parent.Title,
			Description: parent.Description,
			Path:        orchestrator.itemPather().Path(parent.Route().Value()),
		}
	}

	// previous
	if previous := orchestrator.getPrevious(item.Route()); previous != nil {
		navigation.Previous = &viewmodel.NavEntry{
			Title:       previous.Title,
			Description: previous.Description,
			Path:        orchestrator.itemPather().Path(previous.Route().Value()),
		}
	}

	// next
	if next := orchestrator.getNext(item.Route()); next != nil {
		navigation.Next = &viewmodel.NavEntry{
			Title:       next.Title,
			Description: next.Description,
			Path:        orchestrator.itemPather().Path(next.Route().Value()),
		}
	}

	return navigation
}
