// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree/itemtree"
	"github.com/andreaskoch/allmark2/model"
)

func New(logger logger.Logger, repositoryName string) *Index {
	return &Index{
		logger:         logger,
		repositoryName: repositoryName,

		itemList: make([]*model.Item, 0),
		routeMap: make(map[string]*model.Item),
		itemTree: itemtree.New(),
	}
}

type Index struct {
	logger         logger.Logger
	repositoryName string

	// indizes
	itemList []*model.Item
	routeMap map[string]*model.Item // route -> item,
	itemTree *itemtree.ItemTree
}

func (index *Index) IsMatch(route route.Route) (item *model.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[route.Value()]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *Index) IsFileMatch(route route.Route) (*model.File, bool) {

	var parent *model.Item
	parentRoute := &route
	for parentRoute != nil && parentRoute.Level() >= 0 {

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
	if parent == nil {
		index.logger.Warn("No file found for route %q", route)
		return nil, false
	}

	// check if the parent has a file with the supplied route
	if file := parent.GetFile(route); file != nil {
		return file, true
	}

	// file not found
	return nil, false
}

func (index *Index) GetParent(childRoute *route.Route) *model.Item {

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

func (index *Index) Root() *model.Item {
	return index.itemTree.Root()
}

// Get all childs that match the given expression
func (index *Index) GetAllChilds(route *route.Route, expression func(item *model.Item) bool) []*model.Item {

	childs := make([]*model.Item, 0)

	// get all direct childs of the supplied route
	directChilds := index.GetDirectChilds(route)

	for _, child := range directChilds {

		// evaluate expression
		if !expression(child) {
			continue
		}

		// append child
		childs = append(childs, child)

		// recurse
		childs = append(childs, index.GetAllChilds(child.Route(), expression)...)

	}

	return childs
}

func (index *Index) GetDirectChilds(route *route.Route) []*model.Item {
	// get all mathching childs
	childs := index.itemTree.GetChildItems(route)

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
}

// Get a list of all item in this index.
func (index *Index) Items() []*model.Item {
	items := index.itemList
	return items
}

func (index *Index) GetItemByRoute(route *route.Route) (bool, *model.Item) {
	routeValue := route.Value()
	if item, exists := index.routeMap[routeValue]; exists {
		return true, item
	}

	return false, nil
}

func (index *Index) Add(item *model.Item) {

	// abort if item is invalid
	if item == nil {
		index.logger.Warn("Cannot add an invalid item to the index.")
		return
	}

	// make root items a repository
	if item.Route().Level() == 0 {
		item.Type = model.TypeRepository
	}

	index.logger.Debug("Adding item %q to index", item)

	// the the item to the indizes
	index.insertItemToItemList(item)
	index.insertItemToRouteMap(item)
	index.insertItemToTree(item)
}

func (index *Index) insertItemToItemList(item *model.Item) {
	index.itemList = append(index.itemList, item)
}

func (index *Index) insertItemToRouteMap(item *model.Item) {
	itemRoute := item.Route().Value()
	index.routeMap[itemRoute] = item
}

func (index *Index) insertItemToTree(item *model.Item) {
	index.itemTree.InsertItem(item)
}

// sort the items by name
func sortItemsByRoute(item1, item2 *model.Item) bool {

	// ascending by route
	return item1.Route().Value() > item2.Route().Value()

}
