// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"fmt"
)

func New(rootName string, rootValue interface{}) *Tree {
	return &Tree{
		root: newNode(nil, rootName, rootValue),
	}
}

type Tree struct {
	root *Node
}

func (tree *Tree) String() string {
	if tree.Root() == nil {
		return "<empty-tree>"
	}

	return tree.Root().String()
}

func (tree *Tree) Root() *Node {
	return tree.root
}

func (tree *Tree) Insert(path Path, value interface{}) (bool, error) {

	// validate path
	if isValidPath, pathValidationErr := path.IsValid(); !isValidPath {
		return false, pathValidationErr
	}

	// check if the path already exists
	if existingNode := tree.GetNode(path); existingNode != nil {
		existingNode.SetValue(value)
		return true, nil
	}

	// convert components to node
	node := pathToNode(path, value)
	if node == nil {
		return false, fmt.Errorf("Could not convert the path %s into a node.", path)
	}

	// make the new node the root
	if tree.Root() == nil {
		tree.root = node
		return true, nil
	}

	// insert the node
	return tree.Root().Insert(node)
}

func (tree *Tree) Delete(path Path) (bool, error) {

	// validate path
	if isValidPath, pathValidationErr := path.IsValid(); !isValidPath {
		return false, pathValidationErr
	}

	// convert components to node
	node := pathToNode(path, nil)
	if node == nil {
		return false, fmt.Errorf("Could not convert the path %s into a node.", path)
	}

	if tree.Root() == nil {
		return false, fmt.Errorf("Cannot remove the path %s from this tree because the tree is empty.", path)
	}

	// delete the node
	return tree.Root().Delete(node)
}

// Get the node that matches the supplied path.
func (tree *Tree) GetNode(path Path) *Node {

	// validate path
	if isValidPath, _ := path.IsValid(); !isValidPath {
		return nil
	}

	// return nil if there is no root
	if tree.Root() == nil {
		return nil
	}

	// return the root if the path is a root path
	if path.IsRootPath() {
		return tree.Root()
	}

	// skip the root node itself and go for its childs
	for _, child := range tree.Root().Childs() {
		if matchingNode := child.GetNode(path); matchingNode != nil {
			return matchingNode
		}
	}

	// no match
	return nil
}
