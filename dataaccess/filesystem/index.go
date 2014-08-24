// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
)

func newIndex(logger logger.Logger) *Index {
	return &Index{
		logger: logger,

		itemList: make([]*dataaccess.Item, 0),
		routeMap: make(map[string]*dataaccess.Item),
		itemTree: newItemTree(logger),
	}
}

type Index struct {
	logger logger.Logger

	// indizes
	itemList []*dataaccess.Item
	routeMap map[string]*dataaccess.Item // route -> item,
	itemTree *ItemTree
}

func (index *Index) String() string {
	return index.itemTree.String()
}

func (index *Index) IsMatch(route route.Route) (item *dataaccess.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[routeToKey(route)]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *Index) IsFileMatch(route route.Route) (*dataaccess.File, bool) {

	var parent *dataaccess.Item
	parentRoute := route
	for !parentRoute.IsEmpty() && parentRoute.Level() >= 0 {

		parent, _ = index.IsMatch(parentRoute)
		if parent == nil {

			// next level
			newParentRoute, exists := parentRoute.Parent()
			if !exists {
				break
			}

			parentRoute = newParentRoute
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

func (index *Index) GetParent(childRoute route.Route) *dataaccess.Item {

	if childRoute.IsEmpty() {
		return nil
	}

	// abort if the supplied route is already a root
	if childRoute.Level() == 0 {
		return nil
	}

	// get the parent route
	parentRoute, exists := childRoute.Parent()
	if !exists {
		return nil
	}

	item, isMatch := index.IsMatch(parentRoute)
	if !isMatch {
		return nil
	}

	return item
}

func (index *Index) Root() *dataaccess.Item {
	return index.itemTree.Root()
}

// Get all childs that match the given expression
func (index *Index) GetAllChilds(route route.Route, expression func(item *dataaccess.Item) bool) []*dataaccess.Item {

	childs := make([]*dataaccess.Item, 0)

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

func (index *Index) GetDirectChilds(route route.Route) []*dataaccess.Item {
	// get all mathching childs
	childs := index.itemTree.GetChildItems(route)

	// sort the items by ascending by route
	dataaccess.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
}

// Get a list of all item in this index.
func (index *Index) Items() []*dataaccess.Item {
	return index.itemList
}

// Get the number of entries in this index
func (index *Index) Size() int {
	return len(index.itemList)
}

func (index *Index) Add(item *dataaccess.Item) {

	// abort if item is invalid
	if item == nil {
		index.logger.Warn("Cannot add an invalid item to the index.")
		return
	}

	index.logger.Debug("Adding item %q to index", item)

	// the the item to the indizes
	index.addItemToItemList(item)
	index.addItemToRouteMap(item)
	index.addItemToTree(item)
}

func (index *Index) addItemToItemList(item *dataaccess.Item) {
	index.itemList = append(index.itemList, item)
}

func (index *Index) addItemToRouteMap(item *dataaccess.Item) {
	index.routeMap[routeToKey(item.Route())] = item
}

func (index *Index) addItemToTree(item *dataaccess.Item) {
	index.itemTree.Insert(item)
}

func (index *Index) Remove(route route.Route) {

	// locate the item
	item, exists := index.IsMatch(route)
	if !exists {
		index.logger.Warn("The item with the route %q was not found in this index.", route)
		return
	}

	index.logger.Debug("Removing item %q from index", item)

	// the the item to the indizes
	index.removeItemFromItemList(item)
	index.removeItemFromRouteMap(item)
	index.removeItemFromTree(item)
}

func (index *Index) removeItemFromItemList(item *dataaccess.Item) {

	// find the index of the item to remove
	indexToRemove := -1
	for index, child := range index.itemList {
		if item.String() == child.String() {
			indexToRemove = index
			break
		}
	}

	if indexToRemove == -1 {
		// the item was not found
		index.logger.Warn("The item '%s' was not found in the item list.", item)
		return
	}

	index.itemList = append(index.itemList[:indexToRemove], index.itemList[indexToRemove+1:]...)
}

func (index *Index) removeItemFromRouteMap(item *dataaccess.Item) {
	delete(index.routeMap, routeToKey(item.Route()))
}

func (index *Index) removeItemFromTree(item *dataaccess.Item) {
	if _, err := index.itemTree.Delete(item); err != nil {
		index.logger.Error("Unable to delete item '%s' from the item tree. Error: %s", item, err.Error())
	}

}

// sort the items by name
func sortItemsByRoute(item1, item2 *dataaccess.Item) bool {

	// ascending by route
	return item1.Route().Value() > item2.Route().Value()

}
