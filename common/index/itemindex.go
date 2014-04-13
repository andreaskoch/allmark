// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
)

func CreateItemIndex(logger logger.Logger, repositoryName string) *ItemIndex {
	return &ItemIndex{
		logger:         logger,
		repositoryName: repositoryName,
		items:          make(map[route.Route]*model.Item),
	}
}

type ItemIndex struct {
	logger         logger.Logger
	repositoryName string
	items          map[route.Route]*model.Item
}

func (index *ItemIndex) IsMatch(route route.Route) (item *model.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.items[route]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *ItemIndex) IsFileMatch(route route.Route) (*model.File, bool) {

	var parent *model.Item
	parentRoute := &route
	for parentRoute != nil && parentRoute.Level() > 0 {

		parent, _ = index.IsMatch(*parentRoute)
		if parent == nil || parent.IsVirtual() {

			// next level
			parentRoute = parentRoute.Parent()
			continue

		}

		// found a non-virtual parent
		break

	}

	// abort if there is no non-virtual parent
	if parent == nil || parent.IsVirtual() {
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

	if childRoute == nil {
		return nil
	}

	// abort if the supplied route is already a root
	if childRoute.Level() == 0 {
		return nil
	}

	// get the parent route
	parentRoute := childRoute.Parent()
	if parentRoute == nil {
		return nil
	}

	item, isMatch := index.IsMatch(*parentRoute)
	if !isMatch {
		return nil
	}

	return item
}

func (index *ItemIndex) Root() *model.Item {
	root := route.New()

	if _, isMatch := index.IsMatch(*root); !isMatch {

		// create a virtual root
		virtualRoot, err := newVirtualItem(*root)
		if err != nil {
			index.logger.Error("%s", err.Error())
			return nil
		}

		// use the repository name as the title of the root item
		virtualRoot.Title = index.repositoryName

		// write the new virtual root to the index
		index.items[*root] = virtualRoot
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
		index.logger.Debug("%q is at level 0", baseRoute)
		return
	}

	parentRoute := baseRoute.Parent()
	for parentRoute != nil && parentRoute.Level() > 0 {

		if _, exists := index.items[*parentRoute]; !exists {

			index.logger.Debug("Adding virtual item %q to index", parentRoute)

			virtualParentItem, err := newVirtualItem(*parentRoute)
			if err != nil {
				panic(err)
			}

			// add the virtual item to the index
			index.items[*parentRoute] = virtualParentItem

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
	item.Title = route.LastComponentName()

	return item, nil
}

// sort the items by name
func sortItemsByRoute(item1, item2 *model.Item) bool {

	// ascending by route
	return item1.Route().Value() > item2.Route().Value()

}
