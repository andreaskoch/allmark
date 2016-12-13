// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/model"
)

func New(logger logger.Logger) *Index {
	return &Index{
		logger: logger,

		itemList: make([]*model.Item, 0),
		routeMap: make(map[string]*model.Item),
		itemTree: newItemTree(logger),
	}
}

type Index struct {
	logger logger.Logger

	// indizes
	itemList []*model.Item
	routeMap map[string]*model.Item // route -> item,
	itemTree *ItemTree
}

func (index *Index) String() string {
	return index.itemTree.String()
}

func (index *Index) IsMatch(r route.Route) (item *model.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[route.ToKey(r)]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *Index) IsFileMatch(r route.Route) (*model.File, bool) {

	var parent *model.Item
	parentRoute := r
	for parentRoute.Level() >= 0 {

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
		index.logger.Warn("No file found for route %q", r)
		return nil, false
	}

	// check if the parent has a file with the supplied route
	if file := parent.GetFile(r); file != nil {
		return file, true
	}

	// file not found
	return nil, false
}

func (index *Index) GetParent(childRoute route.Route) *model.Item {

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

func (index *Index) Root() *model.Item {
	return index.itemTree.Root()
}

func (index *Index) Size() int {
	return len(index.itemList)
}

// GetAllItems returns all items in the index.
func (index *Index) GetAllItems() []*model.Item {
	items := make([]*model.Item, 0)
	index.itemTree.Walk(func(item *model.Item) {
		items = append(items, item)
	})
	return items
}

// Get all children that match the given expression
func (index *Index) GetAllChildren(route route.Route, expression func(item *model.Item) bool) []*model.Item {

	children := make([]*model.Item, 0)

	// get all direct children of the supplied route
	directChildren := index.GetDirectChildren(route)

	for _, child := range directChildren {

		// evaluate expression
		if !expression(child) {
			continue
		}

		// append child
		children = append(children, child)

		// recurse
		children = append(children, index.GetAllChildren(child.Route(), expression)...)

	}

	// sort the items by ascending by route
	model.SortItemsBy(sortItemsByDate).Sort(children)

	return children
}

func (index *Index) GetDirectChildren(route route.Route) []*model.Item {
	// get all mathching children
	children := index.itemTree.GetChildItems(route)

	// sort the items by ascending by route
	model.SortItemsBy(sortItemsByDate).Sort(children)

	return children
}

func (index *Index) GetLeafes(route route.Route) []*model.Item {

	item := index.itemTree.GetItem(route)
	if item == nil {

		// item not found
		return []*model.Item{}
	}

	// leaf found
	children := index.GetDirectChildren(route)
	if len(children) == 0 {
		return []*model.Item{item}
	}

	// recurse
	leafes := make([]*model.Item, 0)
	for _, child := range children {
		childLeafes := index.GetLeafes(child.Route())
		if len(childLeafes) == 0 {
			continue
		}

		leafes = append(leafes, childLeafes...)
	}

	return leafes
}

func (index *Index) Add(item *model.Item) {

	// abort if item is invalid
	if item == nil {
		index.logger.Warn("Cannot add an invalid item to the index.")
		return
	}

	// the the item to the indizes
	index.itemList = append(index.itemList, item)
	index.routeMap[route.ToKey(item.Route())] = item
	index.itemTree.Insert(item)
}

func (index *Index) Remove(itemRoute route.Route) {
	delete(index.routeMap, route.ToKey(itemRoute))
	index.itemTree.Delete(itemRoute)
}

// sort the models by date and name
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.Before(model2.MetaData.CreationDate)

}
