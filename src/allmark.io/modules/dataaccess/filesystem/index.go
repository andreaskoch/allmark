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
	index.itemTree.Walk(func(item dataaccess.Item) {
		newIndex.Add(item)
	})

	return newIndex
}

func (index *Index) String() string {
	return index.itemTree.String()
}

// IsMatch checks if the specified route can be found in the index.
func (index *Index) IsMatch(r route.Route) (item dataaccess.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[route.ToKey(r)]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

// GetParent returns the parent of the specified route if there is one.
// Otherwise GetParent will return nil.
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

// Size returns the number if items in the index.
func (index *Index) Size() int {
	return len(index.routeMap)
}

// GetAllItems returns a flat list of all items in the index.
func (index *Index) GetAllItems() []dataaccess.Item {
	items := make([]dataaccess.Item, 0)
	index.itemTree.Walk(func(item dataaccess.Item) {
		items = append(items, item)
	})
	return items
}

// Get all children that match the given expression
func (index *Index) GetAllChildren(route route.Route, limitDepth bool, maxDepth int) []dataaccess.Item {

	children := make([]dataaccess.Item, 0)

	if limitDepth {

		// abort if the max depth level has been reached
		if maxDepth == 0 {
			return children
		}

		// count down the max depth
		maxDepth = maxDepth - 1

	}

	// get all direct children of the supplied route
	directChildren := index.GetDirectChildren(route)

	for _, child := range directChildren {

		// append child
		children = append(children, child)

		// recurse
		children = append(children, index.GetAllChildren(child.Route(), limitDepth, maxDepth)...)

	}

	return children
}

func (index *Index) GetLeafes(route route.Route) []dataaccess.Item {

	item := index.itemTree.GetItem(route)
	if item == nil {

		// item not found
		return []dataaccess.Item{}
	}

	// leaf found
	children := index.GetDirectChildren(route)
	if len(children) == 0 {
		return []dataaccess.Item{item}
	}

	// recurse
	leafes := make([]dataaccess.Item, 0)
	for _, child := range children {
		childLeafes := index.GetLeafes(child.Route())
		if len(childLeafes) == 0 {
			continue
		}

		leafes = append(leafes, childLeafes...)
	}

	return leafes
}

func (index *Index) GetSubIndex(subIndexStartRoute route.Route, limitDepth bool, maxDepth int) *Index {

	subindex := newIndex()

	// get the item with the specified route
	root, exists := index.IsMatch(subIndexStartRoute)
	if !exists {

		// return an empty index. There was no item with the given route
		return subindex
	}

	subindex.Add(root)

	for _, child := range index.GetAllChildren(subIndexStartRoute, limitDepth, maxDepth) {
		subindex.Add(child)
	}

	return subindex
}

func (index *Index) GetDirectChildren(route route.Route) []dataaccess.Item {
	// get all mathching children
	children := index.itemTree.GetChildItems(route)

	return children
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

// Remove removes the item with the supplied route from the index.
func (index *Index) Remove(itemRoute route.Route) {
	delete(index.routeMap, route.ToKey(itemRoute))
	index.itemTree.Delete(itemRoute)
}
