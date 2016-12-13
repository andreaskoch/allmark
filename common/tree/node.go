// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"fmt"
)

func newRootNode(name string, value interface{}) *Node {
	return newNode(nil, name, value)
}

func newNode(parent *Node, name string, value interface{}) *Node {

	// create the node
	node := &Node{
		parent: parent,

		name:  name,
		value: value,

		children:       make(map[string]*Node),
		childrenSorted: make([]*Node, 0),
	}

	// add the new node to the parents' children
	if parent != nil {
		parent.Insert(node)
	}

	return node
}

type Node struct {
	name  string
	value interface{}

	parent *Node

	children       map[string]*Node
	childrenSorted []*Node
}

func (node *Node) String() string {
	markdownListIdentifier := "- "
	text := getIndent(node.Level()) + markdownListIdentifier + node.Name()

	for _, child := range node.Children() {
		text += "\n" + child.String()
	}

	return text
}

func (node *Node) Name() string {
	if node.name == "" {
		return "<no-name-assigned>"
	}

	return node.name
}

func (node *Node) Value() interface{} {
	return node.value
}

func (node *Node) SetValue(value interface{}) {
	node.value = value
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) setParent(childNode *Node) {
	node.parent = childNode
}

func (node *Node) Level() int {
	return getNodeLevel(node)
}

func (node *Node) Children() []*Node {
	return node.childrenSorted
}

func (parentNode *Node) Insert(nodeToInsert *Node) (bool, error) {

	if nodeToInsert == nil {
		return false, fmt.Errorf("The supplied node is nil.")
	}

	// make this node the parent for the inserted (child) node
	nodeToInsert.setParent(parentNode)

	// determine the lookup key
	key := nodeToInsert.Name()

	// check if the given node already exists
	existingNode, exists := parentNode.children[key]
	if !exists {

		// insert the child node as a new entry
		parentNode.children[key] = nodeToInsert
		parentNode.childrenSorted = append(parentNode.childrenSorted, nodeToInsert)

		return true, nil // success
	}

	// insert all sub nodes of the supplied child node
	for _, subNode := range nodeToInsert.Children() {
		if childInserted, childInsertErr := existingNode.Insert(subNode); !childInserted {
			return false, childInsertErr
		}
	}

	return true, nil // success
}

func (parentNode *Node) Delete(path Path) (bool, error) {

	if path.IsEmpty() {
		return false, fmt.Errorf("The path is empty.")
	}

	firstComponent := path[0]

	// find a matching child
	matchingChild, matchingChildExists := parentNode.children[firstComponent]
	if !matchingChildExists {
		return false, fmt.Errorf("The node %q was not found.", path)
	}

	// recurse
	if len(path) > 1 {
		return matchingChild.Delete(path[1:])
	}

	// remove the node from the children
	delete(parentNode.children, firstComponent)
	parentNode.childrenSorted = deleteFromSlice(matchingChild, parentNode.childrenSorted)

	return true, nil
}

// Get the node that matches the supplied path.
func (currentNode *Node) GetNode(path Path) *Node {

	if path.IsEmpty() {
		return nil
	}

	// determine the lookup key
	currentNodeName := currentNode.Name()

	// get the first component of the supplied path
	firstPathComponent := path[0]

	// abort if the current node name does not match the first path component
	if currentNodeName != firstPathComponent {
		return nil
	}

	// if we have reached the end of the path we have found a match
	if isLastPathComponent := len(path) == 1; isLastPathComponent {
		return currentNode
	}

	// recurse
	subPath := path[1:]
	for _, childNode := range currentNode.Children() { // todo: make more efficient. lookup instead of iterate.
		if matchingNode := childNode.GetNode(subPath); matchingNode != nil {
			return matchingNode
		}
	}

	// no match found
	return nil
}

// Walk visits the current node, then every child of the current node and then recurses down the children.
func (currentNode *Node) Walk(expression func(node *Node)) {

	// children first
	for _, child := range currentNode.Children() {
		expression(child)
	}

	// recurse
	for _, child := range currentNode.Children() {
		child.Walk(expression)
	}
}

func getNodeLevel(node *Node) int {
	if node == nil {
		panic("Node cannot be nil.")
	}

	if node.Parent() == nil {
		return 0
	}

	return 1 + getNodeLevel(node.Parent())
}

func getIndent(level int) string {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "    "
	}
	return indent
}

func deleteFromSlice(nodeToDelete *Node, nodes []*Node) []*Node {
	newNodes := make([]*Node, 0)

	for _, node := range nodes {
		if node.Name() == nodeToDelete.Name() {
			continue
		}
		newNodes = append(newNodes, node)
	}

	return newNodes
}
