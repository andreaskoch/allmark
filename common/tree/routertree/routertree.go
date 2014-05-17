// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routertree

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree"
	"github.com/andreaskoch/allmark2/common/tree/treeutil"
)

func New() *RouterTree {
	return &RouterTree{
		*tree.New("", nil),
	}
}

type RouterTree struct {
	tree.Tree
}

func (nodeTree *RouterTree) Root() route.Router {
	rootNode := nodeTree.Tree.Root()
	if rootNode == nil {
		return nil
	}

	return nodeToItem(rootNode)
}

func (nodeTree *RouterTree) InsertItem(routerItem route.Router) {

	if routerItem == nil {
		return
	}

	// convert the route to a path
	path := treeutil.RouteToPath(routerItem.Route())

	nodeTree.Tree.Insert(path, routerItem)
}

func (nodeTree *RouterTree) GetRootItem() route.Router {

	// locate the root node
	node := nodeTree.getNode(route.New())
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToItem(node)
}

func (nodeTree *RouterTree) GetItem(route *route.Route) route.Router {

	// locate the node
	node := nodeTree.getNode(route)
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToItem(node)
}

func (nodeTree *RouterTree) GetChildItems(route *route.Route) (childItems []route.Router) {

	node := nodeTree.getNode(route)
	if node == nil {
		return childItems
	}

	for _, childNode := range node.Childs() {
		if routerItem := nodeToItem(childNode); routerItem != nil {
			childItems = append(childItems, routerItem)
		}
	}

	return childItems
}

func (nodeTree *RouterTree) getNode(route *route.Route) *tree.Node {

	if route == nil {
		return nil
	}

	// convert the route to a path
	path := treeutil.RouteToPath(route)

	// locate the node
	node := nodeTree.Tree.GetNode(path)
	if node == nil {
		return nil
	}

	return node
}

func nodeToItem(node *tree.Node) route.Router {
	return node.Value().(route.Router)
}
