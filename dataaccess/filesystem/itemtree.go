// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/tree"
	"github.com/elWyatt/allmark/dataaccess"
	"fmt"
)

func newItemTree() *ItemTree {
	return &ItemTree{
		*tree.Empty(),
	}
}

type ItemTree struct {
	tree.Tree
}

func (itemTree *ItemTree) Root() dataaccess.Item {
	rootNode := itemTree.Tree.Root()
	if rootNode == nil {
		return nil
	}

	return nodeToItem(rootNode)
}

func (itemTree *ItemTree) Insert(item dataaccess.Item) (bool, error) {

	if item == nil {
		return false, fmt.Errorf("Cannot insert item into tree. The supplied item cannot be null.")
	}

	// convert the route to a path
	path := tree.RouteToPath(item.Route())
	return itemTree.Tree.Insert(path, item)
}

func (itemTree *ItemTree) Delete(itemRoute route.Route) (bool, error) {
	return itemTree.Tree.Delete(itemRoute.Components())
}

func (itemTree *ItemTree) GetItem(route route.Route) dataaccess.Item {

	// locate the node
	node := itemTree.getNode(route)
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToItem(node)
}

func (itemTree *ItemTree) GetChildItems(route route.Route) []dataaccess.Item {

	childItems := make([]dataaccess.Item, 0)

	node := itemTree.getNode(route)
	if node == nil {
		return childItems
	}

	for _, childNode := range node.Children() {
		item := nodeToItem(childNode)
		if item == nil {
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
func (itemTree *ItemTree) Walk(expression func(item dataaccess.Item)) {
	itemTree.Tree.Walk(func(node *tree.Node) {
		item := nodeToItem(node)
		if item == nil {
			return
		}

		expression(item)
	})
}

func nodeToItem(node *tree.Node) dataaccess.Item {

	if node.Value() == nil {
		return nil
	}

	return node.Value().(dataaccess.Item)
}
