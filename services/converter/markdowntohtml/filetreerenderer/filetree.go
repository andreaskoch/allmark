// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetreerenderer

import (
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/common/tree"
	"github.com/andreaskoch/allmark/model"
)

type FileNode struct {
	*tree.Node
}

func (f *FileNode) Value() *model.File {
	return nodeToFile(f.Node)
}

func (f *FileNode) Children() []*FileNode {
	children := make([]*FileNode, 0)
	for _, child := range f.Node.Children() {
		children = append(children, &FileNode{child})
	}
	return children
}

func newTree() *FileTree {
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

func (nodeTree *FileTree) Insert(file *model.File) {

	if file == nil {
		return
	}

	// convert the route to a path
	path := tree.RouteToPath(file.Route())

	nodeTree.Tree.Insert(path, file)
}

func (nodeTree *FileTree) GetNode(route route.Route) *FileNode {

	// convert the route to a path
	path := tree.RouteToPath(route)

	// locate the node
	node := nodeTree.Tree.GetNode(path)
	if node == nil {
		return nil
	}

	return &FileNode{node}
}

func nodeToFile(node *tree.Node) *model.File {
	val := node.Value()
	if val == nil {
		return nil
	}

	return val.(*model.File)
}
