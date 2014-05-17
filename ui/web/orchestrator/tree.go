// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree/routertree"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewTreeOrchestrator(itemIndex *index.Index) TreeOrchestrator {
	return TreeOrchestrator{
		itemIndex: itemIndex,
	}
}

type TreeOrchestrator struct {
	itemIndex *index.Index
}

func (orchestrator *TreeOrchestrator) GetTree(pathProvider paths.Pather, routerItems []route.Router) *viewmodel.TreeNode {

	// convert router items to tree
	tree := routertree.New()
	for _, item := range routerItems {
		tree.InsertItem(item)
	}

	// convert tree to viewmodel
	rootItem := tree.GetRootItem()
	if rootItem == nil {
		return nil
	}

	return convertRouterItemToViewModel(*tree, rootItem)
}

func convertRouterItemToViewModel(tree routertree.RouterTree, rootItem route.Router) *viewmodel.TreeNode {

	treeNodeModel := &viewmodel.TreeNode{}

	treeNodeModel.Route = rootItem.Route().Value()
	childs := tree.GetChildItems(rootItem.Route())
	for _, child := range childs {
		childModel := convertRouterItemToViewModel(tree, child)
		if childModel == nil {
			continue
		}

		treeNodeModel.Childs = append(treeNodeModel.Childs, *childModel)
	}

	return treeNodeModel
}
