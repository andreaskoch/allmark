// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/tree"
	"github.com/elWyatt/allmark/model"
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

func (itemTree *ItemTree) Root() *model.Item {
	rootNode := itemTree.Tree.Root()
	if rootNode == nil {
		return nil
	}

	return nodeToItem(rootNode)
}

func (itemTree *ItemTree) Insert(item *model.Item) {

	if item == nil {
		return
	}

	// convert the route to a path
	path := tree.RouteToPath(item.Route())

	if _, err := itemTree.Tree.Insert(path, item); err != nil {
		itemTree.logger.Error("Cannot insert item %q. Error: %s", item.Route(), err.Error())
	}
}

func (itemTree *ItemTree) Delete(itemRoute route.Route) (bool, error) {
	return itemTree.Tree.Delete(itemRoute.Components())
}

func (itemTree *ItemTree) GetItem(route route.Route) *model.Item {

	// locate the node
	node := itemTree.getNode(route)
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToItem(node)
}

func (itemTree *ItemTree) GetChildItems(route route.Route) []*model.Item {

	childItems := make([]*model.Item, 0)

	node := itemTree.getNode(route)
	if node == nil {
		return childItems
	}

	for _, childNode := range node.Children() {
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

// Walk visits every node in the current tree. Starting with the root, every child of the root and then recurses down the children.
func (itemTree *ItemTree) Walk(expression func(item *model.Item)) {
	itemTree.Tree.Walk(func(node *tree.Node) {
		item := nodeToItem(node)
		if item == nil {
			return
		}

		expression(item)
	})
}

func nodeToItem(node *tree.Node) *model.Item {

	if node.Value() == nil {
		return nil
	}

	return node.Value().(*model.Item)
}
