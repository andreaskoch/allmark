// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"testing"
)

func Test_New_ResultNotNil(t *testing.T) {
	// act
	tree := New("root", nil)

	// assert
	if tree == nil {
		t.Errorf("New returned %s, but New should never return nil.", tree)
	}
}

func Test_New_RootNotNil(t *testing.T) {
	// act
	tree := New("Root", nil)

	// assert
	if tree.Root() == nil {
		t.Errorf("The Root of a new tree should not be nil.")
	}
}

func Test_New_NameOfRootIsSet(t *testing.T) {
	// act
	name := "Root"
	tree := New(name, nil)

	// assert
	if tree.Root().Name() != name {
		t.Errorf("The Name of the a new tree should be %q but was %q", name, tree.Root().Name())
	}
}

func Test_Tree_String_EmptyTree(t *testing.T) {
	// arrange
	tree := &Tree{}

	// act
	result := tree.String()

	// assert
	expected := "<empty-tree>"
	if result != expected {
		t.Errorf("The string method of an empty tree should return %q but returned %q instead.", expected, result)
	}
}

func Test_Tree_String_NonEmptyTree(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("root level child 1", "2nd level child"), nil)
	tree.Insert(NewPath("root level child 2", "2nd level child"), nil)

	// act
	result := tree.String()

	// assert
	expected := `- root
    - root level child 1
        - 2nd level child
    - root level child 2
        - 2nd level child`
	if result != expected {
		t.Errorf("The string method of the tree should return %q but returned %q instead.", expected, result)
	}
}

func Test_Tree_Insert_ValidPath_ResultIsTrue(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	result, _ := tree.Insert(NewPath("child 1"), nil)

	// assert
	expected := true
	if result != expected {
		t.Errorf("The Insert method should return %q when a valid path is inserted but returned %q instead.", expected, result)
	}
}

func Test_Tree_Insert_InvalidPath_ResultIsFalse(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	result, _ := tree.Insert(NewPath("child 1 / dasdasd"), nil)

	// assert
	expected := false
	if result != expected {
		t.Errorf("The Insert method should return %q when an invalid path is inserted but returned %q instead.", expected, result)
	}
}

func Test_Tree_Insert_EmptyPath_ResultIsTrue(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	result, _ := tree.Insert(NewPath(), nil)

	// assert
	if result != true {
		t.Errorf("Inserting an empty path should be possible. The root should be overriden but wasn't.")
	}
}

func Test_Tree_Insert_TwoPaths_StructureIsAsExpected(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	tree.Insert(NewPath("child 1"), nil)
	tree.Insert(NewPath("child 2"), nil)

	// assert
	result := tree.String()
	expected := `- root
    - child 1
    - child 2`
	if result != expected {
		t.Errorf("The tree structure should look like this %q but actually looks loke this %q instead.", expected, result)
	}
}

func Test_Tree_Insert_ExistingTree_TreeIsUpdated(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 2 - to be updated"), nil)
	tree.Insert(NewPath("child 2", "sub child 3"), nil)
	tree.Insert(NewPath("child 3"), nil)

	// act
	tree.Insert(NewPath("child 2", "sub child 2 - to be updated", "new child"), nil)

	// assert
	result := tree.String()
	expected := `- root
    - child 1
    - child 2
        - sub child 1
        - sub child 2 - to be updated
            - new child
        - sub child 3
    - child 3`

	if result != expected {
		t.Errorf("The tree structure should look like this %q but actually looks loke this %q instead.", expected, result)
	}
}

func Test_Tree_Delete_ValidPath_ResultIsTrue(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("child 1"), nil)

	// act
	result, _ := tree.Delete(NewPath("child 1"))

	// assert
	expected := true
	if result != expected {
		t.Errorf("The Delete method should return %q when a valid path is deleted but returned %q instead.", expected, result)
	}
}

func Test_Tree_Delete_InvalidPath_ResultIsFalse(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	result, _ := tree.Delete(NewPath("child 1 / dasdasd"))

	// assert
	expected := false
	if result != expected {
		t.Errorf("The Delete method should return %q when an invalid path is deleted but returned %q instead.", expected, result)
	}
}

func Test_Tree_Delete_EmptyPath_ResultIsFalse(t *testing.T) {
	// arrange
	tree := New("root", nil)

	// act
	result, _ := tree.Delete(NewPath())

	// assert
	expected := false
	if result != expected {
		t.Errorf("The Delete method should return %q when an invalid path is deleted but returned %q instead.", expected, result)
	}
}

func Test_Tree_Delete_ThirdChild_StructureIsAsExpected(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("child 1"), nil)
	tree.Insert(NewPath("child 2"), nil)
	tree.Insert(NewPath("child 3"), nil)

	// act
	result, _ := tree.Delete(NewPath("child 3"))

	// assert result
	expected := true
	if result != expected {
		t.Errorf("The Delete method should return %q when an invalid path is deleted but returned %q instead.", expected, result)
	}

	// assert: structure
	stringResult := tree.String()
	expectedStringResult := `- root
    - child 1
    - child 2`
	if stringResult != expectedStringResult {
		t.Errorf("The tree structure should look like this %q but actually looks loke this %q instead.", expectedStringResult, stringResult)
	}
}

func Test_Tree_GetNode_EmptyTree_ResultIsNil(t *testing.T) {
	// arrange
	tree := &Tree{}

	path := NewPath("root", "child 2")
	var expected *Node

	// act
	result := tree.GetNode(path)

	// assert
	if result != expected {
		t.Errorf("Requesting a node from an empty tree should always return %q but returned %q instead.", expected, result)
	}
}

func Test_Tree_GetNode_InvalidPath_ResultIsNil(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 2"), nil)
	tree.Insert(NewPath("child 2", "sub child 3"), nil)
	tree.Insert(NewPath("child 3"), nil)

	path := NewPath("root", "child 2", "sad ads / dasdsadsa")
	var expected *Node

	// act
	result := tree.GetNode(path)

	// assert
	if result != expected {
		t.Errorf("Requesting a node from a tree with an invalid path %q should return %q but returned %q instead.", path, expected, result)
	}
}

func Test_Tree_GetNode_ComplexTree_ExistingPath_ResultIsNotNull(t *testing.T) {
	// arrange
	tree := New("root", nil)
	tree.Insert(NewPath("child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 1"), nil)
	tree.Insert(NewPath("child 2", "sub child 2"), nil)
	tree.Insert(NewPath("child 2", "sub child 3"), nil)
	tree.Insert(NewPath("child 3"), nil)

	path := NewPath("child 2", "sub child 2")

	// act
	result := tree.GetNode(path)

	// assert
	if result == nil {
		t.Errorf("Requesting a node from the tree %s with the path %q should return a node but returned %q instead.", tree, path, result)
	}
}
