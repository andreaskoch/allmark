// Copyright 2013 Andreas Koch. All rights reserved.
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

		childs: make(map[string]*Node),
	}

	// add the new node to the parents' childs
	if parent != nil {
		parent.Insert(node)
	}

	return node
}

type Node struct {
	name  string
	value interface{}

	parent *Node
	childs map[string]*Node
}

func (node *Node) String() string {
	markdownListIdentifier := "- "
	text := getIndent(node.Level()) + markdownListIdentifier + node.Name()

	for _, child := range node.Childs() {
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

func (node *Node) Childs() []*Node {
	childNodes := make([]*Node, 0)
	for _, childNode := range node.childs {
		childNodes = append(childNodes, childNode)
	}
	return childNodes
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
	existingNode, exists := parentNode.childs[key]
	if !exists {

		// insert the child node as a new entry
		parentNode.childs[key] = nodeToInsert

		return true, nil // success
	}

	// insert all sub nodes of the supplied child node
	for _, subNode := range nodeToInsert.Childs() {
		if childInserted, childInsertErr := existingNode.Insert(subNode); !childInserted {
			return false, childInsertErr
		}
	}

	return true, nil // success
}

func (parentNode *Node) Delete(nodeToDelete *Node) (bool, error) {

	if nodeToDelete == nil {
		return false, fmt.Errorf("The supplied node is nil.")
	}

	// determine the lookup key
	key := nodeToDelete.Name()

	// check if the given node already exists
	if _, exists := parentNode.childs[key]; !exists {
		return false, fmt.Errorf("The node %q was not found.", key)
	}

	// remove the node from the childs
	delete(parentNode.childs, key)
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
	for _, childNode := range currentNode.Childs() { // todo: make more efficient. lookup instead of iterate.
		if matchingNode := childNode.GetNode(subPath); matchingNode != nil {
			return matchingNode
		}
	}

	// no match found
	return nil
}

func (node *Node) Walk(walkFunc func(n *Node) bool) {
	if !walkFunc(node) {
		return // stop recursion
	}

	// recurse
	for _, child := range node.Childs() {
		child.Walk(walkFunc)
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
