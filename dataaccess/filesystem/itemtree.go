// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree"
	"github.com/andreaskoch/allmark2/dataaccess"
)

func newItemTree(logger logger.Logger) *ItemTree {
	return &ItemTree{
		*tree.Empty(),
		logger,
	}
}

type ItemTree struct {
	tree.Tree
	logger logger.Logger
}

func (itemTree *ItemTree) Root() *dataaccess.Item {
	rootNode := itemTree.Tree.Root()
	if rootNode == nil {
		return nil
	}

	return nodeToItem(rootNode)
}

func (itemTree *ItemTree) Insert(item *dataaccess.Item) {

	if item == nil {
		return
	}

	// convert the route to a path
	path := tree.RouteToPath(item.Route())

	if _, err := itemTree.Tree.Insert(path, item); err != nil {
		itemTree.logger.Error("Cannot insert item %q. Error: %s", item.Route(), err.Error())
	}
}

func (itemTree *ItemTree) Delete(item *dataaccess.Item) (bool, error) {

	if item == nil {
		return false, fmt.Errorf("The supplied item is empty.")
	}

	itemRoute := item.Route()
	return itemTree.Tree.Delete(itemRoute.Components())
}

func (itemTree *ItemTree) GetItem(route route.Route) *dataaccess.Item {

	// locate the node
	node := itemTree.getNode(route)
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToItem(node)
}

func (itemTree *ItemTree) GetChildItems(route route.Route) []*dataaccess.Item {

	childItems := make([]*dataaccess.Item, 0)

	node := itemTree.getNode(route)
	if node == nil {
		return childItems
	}

	for _, childNode := range node.Childs() {
		item := nodeToItem(childNode)
		if item == nil {
			itemTree.logger.Warn("The item of child node %q is nil.", childNode)
			continue
		}

		childItems = append(childItems, item)
	}

	return childItems
}

func (itemTree *ItemTree) getNode(route route.Route) *tree.Node {

	// convert the route to a path
	path := tree.RouteToPath(route)

	// locate the node
	node := itemTree.Tree.GetNode(path)
	if node == nil {
		return nil
	}

	return node
}

func nodeToItem(node *tree.Node) *dataaccess.Item {

	if node.Value() == nil {
		return nil
	}

	return node.Value().(*dataaccess.Item)
}
