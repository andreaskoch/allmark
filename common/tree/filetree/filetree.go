// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetree

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree"
	"github.com/andreaskoch/allmark2/common/tree/treeutil"
	"github.com/andreaskoch/allmark2/model"
)

func New() *FileTree {
	return &FileTree{
		*tree.New("", nil),
	}
}

type FileTree struct {
	tree.Tree
}

func (nodeTree *FileTree) Root() *model.File {
	rootNode := nodeTree.Tree.Root()
	if rootNode == nil {
		return nil
	}

	return nodeToFile(rootNode)
}

func (nodeTree *FileTree) InsertFile(file *model.File) {

	if file == nil {
		return
	}

	// convert the route to a path
	path := treeutil.RouteToPath(file.Route())

	nodeTree.Tree.Insert(path, file)
}

func (nodeTree *FileTree) GetFile(route *route.Route) *model.File {

	// locate the node
	node := nodeTree.getNode(route)
	if node == nil {
		return nil
	}

	// cast the value
	return nodeToFile(node)
}

func (nodeTree *FileTree) GetChildFiles(route *route.Route) []*model.File {

	childFiles := make([]*model.File, 0)

	node := nodeTree.getNode(route)
	if node == nil {
		return childFiles
	}

	for _, childNode := range node.Childs() {
		if file := nodeToFile(childNode); file != nil {
			childFiles = append(childFiles, file)
		}
	}

	return childFiles
}

func (nodeTree *FileTree) getNode(route *route.Route) *tree.Node {

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

func nodeToFile(node *tree.Node) *model.File {
	return node.Value().(*model.File)
}
