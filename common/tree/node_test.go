// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"strings"
	"testing"
)

func Test_newNode_ResultNotNil(t *testing.T) {
	// arrange
	var parent *Node
	name := ""

	// act
	result := newNode(parent, name, nil)

	// assert
	if result == nil {
		t.Errorf("The result of newNode should never be nil.")
	}
}

func Test_Node_String_SingleNode_ResultMarkdownListStyle(t *testing.T) {
	// arrange
	name := "Some name"
	node := newNode(nil, name, nil)

	// act
	result := node.String()

	// assert
	expected := "- " + name
	if result != expected {
		t.Errorf("The String method of the node should return %q but returned %q instead.", expected, result)
	}
}

func Test_Node_String_MultipleNodes_ResultIsIndentedMarkdownList(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	firstLevelNode := newNode(rootNode, "first level node", nil)
	newNode(firstLevelNode, "first level child node 1", nil)
	newNode(firstLevelNode, "first level child node 2", nil)

	// act
	result := rootNode.String()

	// assert
	expected := `- root
    - first level node
        - first level child node 1
        - first level child node 2`
	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf(`The String method of the node should return %q but returned %q instead.`, expected, result)
	}
}

func Test_Node_Name_EmptyNode_ResultIsSpecifiedNodeName(t *testing.T) {
	// arrange
	node := &Node{}

	// act
	result := node.Name()

	// assert
	expected := "<no-name-assigned>"
	if result != expected {
		t.Errorf("The name of an empty node should be %q but was %q instead.", expected, result)
	}
}

func Test_Node_Name_NonEmptyNode_ResultIsSpecifiedNodeName(t *testing.T) {
	// arrange
	name := "Some name"
	node := newNode(nil, name, nil)

	// act
	result := node.Name()

	// assert
	if result != name {
		t.Errorf("The name of the node should be %q but was %q.", name, result)
	}
}

func Test_Node_Value_IsNil_ResultIsNil(t *testing.T) {
	// arrange
	var expected interface{}
	node := newRootNode("root", expected)

	// act
	result := node.Value()

	// assert
	if result != expected {
		t.Errorf("The node value should be %q but was %q instead.", expected, result)
	}
}

func Test_Node_Value_IsNotNull_ResultAsSpecified(t *testing.T) {
	// arrange
	expected := new(interface{})
	node := newRootNode("root", expected)

	// act
	result := node.Value()

	// assert
	if result != expected {
		t.Errorf("The node value should be %q but was %q instead.", expected, result)
	}
}

func Test_Node_Parent_NoParentSpecified_ResultIsNil(t *testing.T) {
	// arrange
	node := newNode(nil, "Some name", nil)

	// act
	result := node.Parent()

	// assert
	if result != nil {
		t.Errorf("The parent of the node should be %q but was %q.", nil, result)
	}
}

func Test_Node_Parent_ParentSpecified_ResultIsSpecifiedParent(t *testing.T) {
	// arrange
	parent := newNode(nil, "root", nil)
	node := newNode(parent, "Some name", nil)

	// act
	result := node.Parent()

	// assert
	if result != parent {
		t.Errorf("The parent of the node should be %q but was %q.", parent, result)
	}
}

func Test_Node_GetNode_ThreeLevelNode_EmptyPath_ResultIsNil(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	nodeLevel1 := newNode(rootNode, "level 1", nil)
	newNode(nodeLevel1, "level 2", nil)

	path := NewPath()
	var expectedResult *Node

	// act
	result := rootNode.GetNode(path)

	// assert
	if result != expectedResult {
		t.Errorf("The GetNode method executed on the node %q should return %s for the path %q but returned %s instead.", rootNode, expectedResult, path, result)
	}
}

func Test_Node_GetNode_ThreeLevelNode_FirstNodeIsRequested_ResultIsFirstNode(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	nodeLevel1 := newNode(rootNode, "level 1", nil)
	newNode(nodeLevel1, "level 2", nil)

	path := NewPath("root")
	expectedResult := rootNode

	// act
	result := rootNode.GetNode(path)

	// assert
	if result != expectedResult {
		t.Errorf("The GetNode method executed on the node %q should return %s for the path %q but returned %s instead.", rootNode, expectedResult.Name(), path, result.Name())
	}
}

func Test_Node_GetNode_ThreeLevelNode_SecondlastNodeIsRequested_ResultIsSecondlastNode(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	nodeLevel1 := newNode(rootNode, "level 1", nil)
	newNode(nodeLevel1, "level 2", nil)

	path := NewPath("root", "level 1")
	expectedResult := nodeLevel1

	// act
	result := rootNode.GetNode(path)

	// assert
	if result != expectedResult {
		t.Errorf("The GetNode method executed on the node %q should return %s for the path %q but returned %s instead.", rootNode, expectedResult.Name(), path, result.Name())
	}
}

func Test_Node_GetNode_ThreeLevelNode_LastNodeIsRequested_ResultIsLastNode(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	nodeLevel1 := newNode(rootNode, "level 1", nil)
	nodeLevel2 := newNode(nodeLevel1, "level 2", nil)

	path := NewPath("root", "level 1", "level 2")
	expectedResult := nodeLevel2

	// act
	result := rootNode.GetNode(path)

	// assert
	if result != expectedResult {
		t.Errorf("The GetNode method executed on the node %q should return %s for the path %q but returned %s instead.", rootNode, expectedResult.Name(), path, result.Name())
	}
}

func Test_Node_GetNode_ComplextNode_ComplextPath_ResultIsCorrect(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	newNode(rootNode, "child1", nil)
	newNode(rootNode, "child2", nil)
	newNode(rootNode, "child3", nil)

	nodeLevel1 := newNode(rootNode, "level 1", nil)
	newNode(nodeLevel1, "level 1 child 1", nil)
	newNode(nodeLevel1, "level 1 child 2", nil)
	nodeLevel1Child3 := newNode(nodeLevel1, "level 1 child 3", nil)
	targetNode := newNode(nodeLevel1Child3, "another child", nil)

	newNode(nodeLevel1, "level 2", nil)

	path := NewPath("root", "level 1", "level 1 child 3", "another child")
	expectedResult := targetNode

	// act
	result := rootNode.GetNode(path)

	// assert
	if result != expectedResult {
		t.Errorf("The GetNode method executed on the node %q should return %s for the path %q but returned %s instead.", rootNode, expectedResult.Name(), path, result.Name())
	}
}

func Test_getNodeLevel_RootNode_ResultIsZero(t *testing.T) {
	// arrange
	node := newNode(nil, "node", nil)

	// act
	result := getNodeLevel(node)

	// assert
	expected := 0
	if result != expected {
		t.Errorf("The level of a root node should be %v but was %v.", expected, result)
	}
}

func Test_getNodeLevel_OneParent_ResultIsOne(t *testing.T) {
	// arrange
	firstParent := newNode(nil, "first parent", nil)
	node := newNode(firstParent, "node", nil)

	// act
	result := getNodeLevel(node)

	// assert
	expected := 1
	if result != expected {
		t.Errorf("The level of a root node should be %v but was %v.", expected, result)
	}
}

func Test_getNodeLevel_TwoParents_ResultIsTwo(t *testing.T) {
	// arrange
	firstParent := newNode(nil, "first parent", nil)
	secondParent := newNode(firstParent, "second parent", nil)
	node := newNode(secondParent, "node", nil)

	// act
	result := getNodeLevel(node)

	// assert
	expected := 2
	if result != expected {
		t.Errorf("The level of a root node should be %v but was %v.", expected, result)
	}
}

func Test_getNodeLevel_ThreeParents_ResultIsThree(t *testing.T) {
	// arrange
	firstParent := newNode(nil, "first parent", nil)
	secondParent := newNode(firstParent, "second parent", nil)
	thirdParent := newNode(secondParent, "third parent", nil)
	node := newNode(thirdParent, "node", nil)

	// act
	result := getNodeLevel(node)

	// assert
	expected := 3
	if result != expected {
		t.Errorf("The level of a root node should be %v but was %v.", expected, result)
	}
}

func Test_Node_Delete(t *testing.T) {
	// arrange
	rootNode := newNode(nil, "root", nil)
	firstLevelNode := newNode(rootNode, "first level node", nil)
	secondLevelNode := newNode(firstLevelNode, "first level child node 1", nil)
	newNode(secondLevelNode, "Third level node 1", nil)
	newNode(firstLevelNode, "first level child node 2", nil)

	// act
	rootNode.Delete(Path{"first level node", "first level child node 1", "Third level node 1"})
	result := rootNode.String()

	// assert
	expected := `- root
    - first level node
        - first level child node 1
        - first level child node 2`
	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf(`The String method of the node should return %q but returned %q instead.`, expected, result)
	}
}
