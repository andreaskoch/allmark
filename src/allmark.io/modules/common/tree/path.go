// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"fmt"
	"strings"
)

func NewPath(component ...string) Path {
	return component
}

type Path []string

func (path Path) String() string {
	if path.IsEmpty() {
		return "<empty-path>"
	}

	return strings.Join(path, " > ")
}

func (path Path) IsEmpty() bool {
	return len(path) == 0
}

func (path Path) IsRootPath() bool {
	// contains only one component and this component is empty
	return len(path) == 1 && path[0] == ""
}

// Checks whether all components of the current path are valid (not empty, and don't contain slashes).
func (path Path) IsValid() (bool, error) {

	// abort if its a root path
	if path.IsRootPath() || path.IsEmpty() {
		return true, nil
	}

	for _, component := range path {
		if isValidComponent, errMessage := isValidPathComponent(component); !isValidComponent {
			return false, errMessage
		}
	}

	return true, nil
}

func isValidPathComponent(component string) (bool, error) {

	// empty
	if isEmpty := component == ""; isEmpty {
		return false, fmt.Errorf("Path component cannot be empty")
	}

	// forward slashes
	if containsForwardSlash := strings.Contains(component, "/"); containsForwardSlash {
		return false, fmt.Errorf("Path components cannot contain forward slashes (/).")
	}

	// backslash
	if containsBackslash := strings.Contains(component, `\`); containsBackslash {
		return false, fmt.Errorf(`Path components cannot contain backslashes (\).`)
	}

	return true, nil
}

func pathToNode(path Path, value interface{}) *Node {

	length := len(path)
	if length == 0 {
		return nil
	}

	var firstNode *Node

	var parent *Node
	for index, component := range path {

		var nodeValue interface{}
		if isLastComponent := index+1 == length; isLastComponent {
			nodeValue = value
		}

		node := newNode(parent, component, nodeValue)

		// capture the first node
		if isFirstNode := parent == nil; isFirstNode {
			firstNode = node
		}

		// make this node the parent for the next node
		parent = node
	}

	return firstNode
}
