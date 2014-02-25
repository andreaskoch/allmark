// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
)

func CreateItemIndex(logger logger.Logger) *ItemIndex {
	return &ItemIndex{
		logger: logger,
		items:  make(map[route.Route]*model.Item),
	}
}

type ItemIndex struct {
	logger logger.Logger
	items  map[route.Route]*model.Item
}

func (index *ItemIndex) IsMatch(route route.Route) (item *model.Item, isMatch bool) {
	item, isMatch = index.items[route]
	return
}

func (index *ItemIndex) IsFileMatch(route route.Route) (*model.File, bool) {

	fmt.Printf("Checking if %q is a file match\n", route)

	// skip all virtual parents
	parent := index.GetParent(&route)
	for parent != nil && parent.IsVirtual() {
		fmt.Printf("Parent is %q (IsVirtual: %s)\n", parent.Route(), parent.IsVirtual())
		parent = index.GetParent(parent.Route())
	}

	// abort if there is no non-virtual parent
	if parent == nil {
		return nil, false
	}

	// check if the parent has a file with the supplied route
	if file := parent.GetFile(route); file != nil {
		return file, true
	}

	// file not found
	return nil, false
}

func (index *ItemIndex) GetParent(childRoute *route.Route) *model.Item {

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
		fmt.Printf("%q is a parent of %q\n", parentRoute, childRoute)
		return parentItem
	}

	// check if there is a parent
	parentRoute := childRoute.Parent()
	if parentRoute == nil {
		return nil // we are already at the root
	}

	// no parent item found - create a virtual parent
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

func (index *ItemIndex) GetChilds(route *route.Route) []*model.Item {

	// routeLevel := route.Level()
	childs := make([]*model.Item, 0)

	for itemRoute, item := range index.items {

		// skip all items which are not a child
		if !itemRoute.IsChildOf(route) {
			continue
		}

		childs = append(childs, item)
	}

	// insert virtual items
	childs = index.FillGapsWithVirtualItems(route, childs)

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
}

func (index *ItemIndex) FillGapsWithVirtualItems(baseRoute *route.Route, items []*model.Item) []*model.Item {
	baseRouteLevel := baseRoute.Level()

	itemsByRouteLevel := make(map[int]*model.Item)
	for _, item := range items {
		routeLevel := item.Route().Level()

		// store the route by its level
		itemsByRouteLevel[routeLevel] = item
	}

	// locate the gabs and fill them
	newItems := items
	for level, item := range itemsByRouteLevel {

		// skip the base route level
		if level == baseRouteLevel {
			continue
		}

		// there cannot be a parent for level zero
		if level == 0 {
			continue
		}

		// check if there is a parent for the current item/level
		if _, hasParentInList := itemsByRouteLevel[level-1]; !hasParentInList {

			// get a (virtual) parent from the index
			if newParent := index.GetParent(item.Route()); newParent != nil {
				newItems = append(newItems, newParent)
			}

		}
	}

	return newItems
}

func (index *ItemIndex) Routes() []route.Route {
	routes := make([]route.Route, 0)
	for route, _ := range index.items {
		routes = append(routes, route)
	}
	return routes
}

func (index *ItemIndex) Items() []*model.Item {
	items := make([]*model.Item, 0, len(index.items))

	for _, item := range index.items {
		items = append(items, item)
	}

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(items)

	return items
}

// Get the maxium level of all routes in this index (default: 0)
func (index *ItemIndex) MaxLevel() int {

	maxLevel := 0

	for _, item := range index.items {
		itemLevel := item.Route().Level()
		if itemLevel > maxLevel {
			maxLevel = itemLevel
		}
	}

	return maxLevel
}

func (index *ItemIndex) Add(item *model.Item) {
	index.logger.Debug("Adding item %q to index", item)
	index.items[*item.Route()] = item
}

// sort the items by date and name
func sortItemsByRoute(item1, item2 *model.Item) bool {

	// ascending by route
	return item1.Route().Value() < item2.Route().Value()

}
