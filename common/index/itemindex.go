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

		itemList: make([]*model.Item, 0),
		routeMap: make(map[string]*model.Item),
		itemTree: NewItemTree(),
	}
}

type ItemIndex struct {
	logger         logger.Logger
	repositoryName string

	// indizes
	itemList []*model.Item
	routeMap map[string]*model.Item // route -> item,
	itemTree *ItemTree
}

func (index *ItemIndex) IsMatch(route route.Route) (item *model.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.routeMap[route.Value()]; isMatch {
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
		index.addVirtualItem(*root)
	}

	return index.itemTree.Root()
}

// Get all childs that match the given expression
func (index *ItemIndex) GetAllChilds(route *route.Route, expression func(item *model.Item) bool) []*model.Item {

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

func (index *ItemIndex) GetDirectChilds(route *route.Route) []*model.Item {
	// get all mathching childs
	childs := index.itemTree.GetChildItems(route)

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
}

// Get a list of all item in this index.
func (index *ItemIndex) Items() []*model.Item {
	items := index.itemList
	return items
}

func (index *ItemIndex) GetItemByRoute(route *route.Route) (bool, *model.Item) {
	routeValue := route.Value()
	if item, exists := index.routeMap[routeValue]; exists {
		return true, item
	}

	return false, nil
}

func (index *ItemIndex) Add(item *model.Item) {

	// abort if item is invalid
	if item == nil {
		index.logger.Warn("Cannot add an invalid item to the index.")
		return
	}

	index.logger.Debug("Adding item %q to index", item)

	// the the item to the indizes
	index.insertItemToItemList(item)
	index.insertItemToRouteMap(item)
	index.insertItemToTree(item)

	// insert virtual items if required
	index.addVirtualItem(*item.Route())
}

func (index *ItemIndex) insertItemToItemList(item *model.Item) {
	index.itemList = append(index.itemList, item)
}

func (index *ItemIndex) insertItemToRouteMap(item *model.Item) {
	itemRoute := item.Route().Value()
	index.routeMap[itemRoute] = item
}

func (index *ItemIndex) insertItemToTree(item *model.Item) {
	index.itemTree.InsertItem(item)
}

func (index *ItemIndex) addVirtualItem(baseRoute route.Route) {

	// validate the input
	if baseRoute.Level() == 0 {
		index.logger.Debug("%q is at level 0", baseRoute)
		return
	}

	parentRoute := baseRoute.Parent()
	if _, exists := index.IsMatch(*parentRoute); !exists {

		// create a new virtual item
		index.logger.Debug("Adding virtual item %q to index", parentRoute)
		virtualParentItem, err := newVirtualItem(*parentRoute)
		if err != nil {
			panic(err)
		}

		// use the repository name as the title if the item it the root
		if isRoot := parentRoute.Level() == 0; isRoot {
			virtualParentItem.Title = index.repositoryName
		}

		// add the virtual item to the index
		index.Add(virtualParentItem)
	}
}

func newVirtualItem(route route.Route) (*model.Item, error) {

	// determine the item type
	itemType := model.TypeDocument
	if route.Level() == 0 {

		// root item get the type "repository"
		itemType = model.TypeRepository

	}

	// create a virtual item
	item, err := model.NewVirtualItem(&route, itemType)
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
