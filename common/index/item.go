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

	// locate first and second level childs
	childs := make([]*model.Item, 0)

	for itemRoute, item := range index.items {

		// skip all items which are not a child
		if !itemRoute.IsChildOf(route) {
			continue
		}

		childs = append(childs, item)
	}

	return childs
}

func (index *ItemIndex) Routes() []route.Route {
	routes := make([]route.Route, 0)
	for route, _ := range index.items {
		routes = append(routes, route)
	}
	return routes
}

func (index *ItemIndex) Add(item *model.Item) {
	index.logger.Debug("Adding item %q to index", item)
	index.items[*item.Route()] = item
}
