// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
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

	// check for a direct match
	if item, isMatch = index.items[route]; isMatch {
		return item, isMatch
	}

	// the route has childs we can create a virtual item for it
	if hasChilds := len(index.GetAllChilds(&route)) > 0; hasChilds {

		// if there is an indirect match we can return a virtual item
		virtualItem, err := newVirtualItem(route)
		if err != nil {
			index.logger.Error("Could not create a virtual item for route %q. Error: %s", route, err)
			return nil, false
		}

		return virtualItem, true

	}

	// no match
	return nil, false
}

func (index *ItemIndex) IsFileMatch(route route.Route) (*model.File, bool) {

	// skip all virtual parents
	parent := index.GetParent(&route)
	for parent != nil && parent.IsVirtual() {
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
		return parentItem
	}

	// check if there is a parent
	parentRoute := childRoute.Parent()
	if parentRoute == nil {
		return nil // we are already at the root
	}

	// no parent item found - create a virtual parent
	virtualParent, err := newVirtualItem(*parentRoute)
	if err != nil {

		// error while creating a virtual parent
		index.logger.Warn("Unable to create a virtual parent for the route %q. Error: %s", parentRoute, err)
		return nil

	}

	// store the virtual parent in the index
	index.items[*parentRoute] = virtualParent

	// return the virtual parent
	return virtualParent
}

func (index *ItemIndex) Root() *model.Item {
	root, err := route.New()
	if err != nil {
		return nil
	}

	return index.items[*root]
}

func (index *ItemIndex) GetAllChilds(route *route.Route) []*model.Item {
	return index.getChilds(route, true)
}

func (index *ItemIndex) GetChilds(route *route.Route) []*model.Item {
	return index.getChilds(route, false)
}

func (index *ItemIndex) getChilds(route *route.Route, recurse bool) []*model.Item {

	routeLevel := route.Level()
	nextLevel := routeLevel + 1

	// routeLevel := route.Level()
	childs := make([]*model.Item, 0)

	for childRoute, child := range index.items {

		// skip all deeper-level childs if recursion is disabled
		if !recurse && childRoute.Level() != nextLevel {
			continue
		}

		// skip all items which are not a child
		if !childRoute.IsChildOf(route) {
			continue
		}

		childs = append(childs, child)
	}

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
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

	// the the item to the index
	itemRoute := *item.Route()
	index.items[itemRoute] = item

	// insert virtual items if required
	index.fillGapsWithVirtualItems(itemRoute)
}

func (index *ItemIndex) fillGapsWithVirtualItems(baseRoute route.Route) {

	// validate the input
	if baseRoute.Level() == 0 {
		return
	}

	parentRoute := baseRoute.Parent()
	for parentRoute != nil && parentRoute.Level() > 0 {

		if _, exists := index.items[*parentRoute]; !exists {

			if virtualParentItem, err := newVirtualItem(*parentRoute); err != nil {

				// add the virtual item to the index
				index.items[*parentRoute] = virtualParentItem

			}
		}

		// move up
		parentRoute = parentRoute.Parent()

	}
}

func newVirtualItem(route route.Route) (*model.Item, error) {

	// create a virtual item
	item, err := model.NewVirtualItem(&route)
	if err != nil {
		return nil, err
	}

	// set the item title
	item.Title = route.FolderName()

	return item, nil
}

// sort the items by date and name
func sortItemsByRoute(item1, item2 *model.Item) bool {

	// ascending by route
	return item1.Route().Value() < item2.Route().Value()

}
