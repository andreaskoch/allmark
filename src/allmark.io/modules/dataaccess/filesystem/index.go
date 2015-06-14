// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
	"fmt"
)

func newIndex() *Index {
	return &Index{
		routeMap: make(map[string]dataaccess.Item),
		itemTree: newItemTree(),
	}
}

type Index struct {

	// indizes
	routeMap map[string]dataaccess.Item // route -> item,
	itemTree *ItemTree
}

// Copy creates a copy of the current index
func (index *Index) Copy() *Index {
	newIndex := newIndex()
	for _, existingItem := range index.GetAllItems() {
		newIndex.Add(existingItem)
	}

	return newIndex
}

func (index *Index) IsMatch(r route.Route) (item dataaccess.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[route.ToKey(r)]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *Index) GetParent(childRoute route.Route) dataaccess.Item {

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

func (index *Index) Size() int {
	return len(index.routeMap)
}

// Get all items
func (index *Index) GetAllItems() []dataaccess.Item {
	items := make([]dataaccess.Item, 0)
	index.itemTree.Walk(func(item dataaccess.Item) {
		items = append(items, item)
	})
	return items
}

// Get all childs that match the given expression
func (index *Index) GetAllChilds(route route.Route, expression func(item dataaccess.Item) bool) []dataaccess.Item {

	childs := make([]dataaccess.Item, 0)

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

func (index *Index) GetDirectChilds(route route.Route) []dataaccess.Item {
	// get all mathching childs
	childs := index.itemTree.GetChildItems(route)

	return childs
}

func (index *Index) GetLeafes(route route.Route) []dataaccess.Item {

	item := index.itemTree.GetItem(route)
	if item == nil {

		// item not found
		return []dataaccess.Item{}
	}

	// leaf found
	childs := index.GetDirectChilds(route)
	if len(childs) == 0 {
		return []dataaccess.Item{item}
	}

	// recurse
	leafes := make([]dataaccess.Item, 0)
	for _, child := range childs {
		childLeafes := index.GetLeafes(child.Route())
		if len(childLeafes) == 0 {
			continue
		}

		leafes = append(leafes, childLeafes...)
	}

	return leafes
}

func (index *Index) Add(item dataaccess.Item) (bool, error) {

	// abort if item is invalid
	if item == nil {
		return false, fmt.Errorf("Cannot add nil item to index.")
	}

	// the the item to the indizes
	index.routeMap[route.ToKey(item.Route())] = item
	return index.itemTree.Insert(item)
}

func (index *Index) Remove(itemRoute route.Route) {
	delete(index.routeMap, route.ToKey(itemRoute))
	index.itemTree.Delete(itemRoute)
}
