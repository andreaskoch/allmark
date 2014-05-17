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

func (orchestrator *TreeOrchestrator) GetTree(pathProvider paths.Pather, routerItems []route.Router) viewmodel.TreeNode {

	// convert router items to tree
	tree := routertree.New()
	for _, item := range routerItems {
		tree.InsertItem(item)
	}

	// convert tree to viewmodel
	return convert(*tree)
}

func convert(tree routertree.RouterTree) viewmodel.TreeNode {

	viewModel := viewmodel.TreeNode{}

	tree.WalkItems(tree.Root(), func(router *route.Router) bool {
		viewModel.Route = router.Route()
		return true
	})

	return viewModel
}
