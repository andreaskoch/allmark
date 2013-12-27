// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
)

func New(logger logger.Logger) *Index {
	return &Index{
		logger: logger,
		items:  make(map[route.Route]*model.Item),
	}
}

type Index struct {
	logger logger.Logger
	items  map[route.Route]*model.Item
}

func (index *Index) IsMatch(route route.Route) (item *model.Item, isMatch bool) {
	item, isMatch = index.items[route]
	return
}

func (index *Index) GetParent(item *model.Item) *model.Item {
	childRoute := item.Route()

	// already at the root
	if childRoute.Level() == 0 {
		return nil
	}

	// locate the parent item
	for parentRoute, parentItem := range index.items {
		if !parentRoute.IsParentOf(childRoute) {
			continue
		}

		// return the parent
		return parentItem
	}

	// no parent item found - create a virtual parent
	parentRoute := childRoute.Parent()
	virtualParent, err := newVirtualItem(parentRoute)
	if err != nil {

		// error while creating a virtual parent
		index.logger.Warn("Unable to create a virtual parent for the route %q. Error: %s", parentRoute, err)
		return nil

	}

	// return the virtual parent
	return virtualParent
}

func newVirtualItem(route *route.Route) (*model.Item, error) {

	// create a virtual item
	item, err := model.NewVirtualItem(route)
	if err != nil {
		return nil, err
	}

	// set the item title
	item.Title = route.FolderName()

	return item, nil
}

func (index *Index) GetChilds(item *model.Item) []*model.Item {
	route := item.Route()
	routeLevel := route.Level()

	// locate first and second level childs
	firstLevelChilds := make([]*model.Item, 0)
	deeperLevelChilds := make([]*model.Item, 0)

	for itemRoute, item := range index.items {

		// skip all items which are not a child
		if !itemRoute.IsChildOf(route) {
			continue
		}

		// a first level child is one that is one level deeper than the parent
		isFirstLevelChild := itemRoute.Level() == routeLevel+1

		// add first and second level childs to their respective lists
		if isFirstLevelChild {
			firstLevelChilds = append(firstLevelChilds, item)
		} else {
			deeperLevelChilds = append(deeperLevelChilds, item)
		}
	}

	// just return the first level childs if there are no deeper level childs
	if len(deeperLevelChilds) == 0 {
		return firstLevelChilds
	}

	// create virtual items from the deeper level item if there is no first-level child
	secondLevelChildMap := make(map[string]*model.Item)
	for _, deeperLevelChild := range deeperLevelChilds {

		// determine child route
		secondLevelChildRoute, _ := deeperLevelChild.Route().SubRoute(routeLevel + 1)
		if secondLevelChildRoute == nil {
			continue
		}

		// skip existing entries
		if _, exists := secondLevelChildMap[secondLevelChildRoute.Value()]; exists {
			continue
		}

		// create virtual child item
		virtualChild, err := newVirtualItem(secondLevelChildRoute)
		if err != nil {
			index.logger.Warn("Unable to create a virtual child for the route %q. Error: %s", secondLevelChildRoute, err)
			continue
		}

		secondLevelChildMap[secondLevelChildRoute.Value()] = virtualChild

	}

	secondLevelChilds := make([]*model.Item, 0)
	for _, secondLevelChild := range secondLevelChildMap {
		secondLevelChilds = append(secondLevelChilds, secondLevelChild)
	}

	return secondLevelChilds
}

func (index *Index) Routes() []route.Route {
	routes := make([]route.Route, 0)
	for route, _ := range index.items {
		routes = append(routes, route)
	}
	return routes
}

func (index *Index) Add(item *model.Item) {
	index.logger.Debug("Adding item %q to index", item)
	index.items[*item.Route()] = item
}
