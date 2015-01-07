// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"testing"
)

func Test_Path_String_EmptyPath_ResultIsDefaultString(t *testing.T) {
	// arrange
	path := Path{}

	// act
	result := path.String()

	// assert
	expected := "<empty-path>"
	if result != expected {
		t.Errorf("The String method of an empty path should return %q but returned %q instead.", expected, result)
	}
}

func Test_Path_String_NonEmptyPath_ResultIsCompoentsConnectedByArrows(t *testing.T) {
	// arrange
	path := Path{"root", "level 1", "level 2"}

	// act
	result := path.String()

	// assert
	expected := "root > level 1 > level 2"
	if result != expected {
		t.Errorf("The String method of the Path should return %q but returned %q instead.", expected, result)
	}
}

func Test_Path_IsEmpty_EmptyPath_ResultIsTrue(t *testing.T) {
	// arrange
	path := Path{}

	// act
	result := path.IsEmpty()

	// assert
	expected := true
	if result != expected {
		t.Errorf("The IsEmpty method of an empty path should return %v but returned %v instead.", expected, result)
	}
}

func Test_Path_IsEmpty_NonEmptyPath_ResultIsFalse(t *testing.T) {
	// arrange
	path := Path{"root", "level 1"}

	// act
	result := path.IsEmpty()

	// assert
	expected := false
	if result != expected {
		t.Errorf("The IsEmpty method of a non-empty path should return %v but returned %v instead.", expected, result)
	}
}

func Test_Path_IsValid_AllValidComponents_ResultIsTrue(t *testing.T) {
	// arrange
	path := Path{"root", "level 1", "öäöölasdalsöä", " dsad&42323", "5666%", "_dasdsadas_dsaasd-dasdas?&"}

	// act
	result, _ := path.IsValid()

	// assert
	expected := true
	if result != expected {
		t.Errorf("The IsValid method of the path %q should return %v but returned %v instead.", path, expected, result)
	}
}

func Test_Path_IsValid_ASingleInvalidComponent_ResultIsFalse(t *testing.T) {
	// arrange
	path := Path{"root", "level 1", "öäöölasdalsöä", " dsad / 42323", "5666%", "_dasdsadas_dsaasd-dasdas?&"}

	// act
	result, _ := path.IsValid()

	// assert
	expected := false
	if result != expected {
		t.Errorf("The IsValid method of the path %q should return %v but returned %v instead.", path, expected, result)
	}
}

func Test_pathToNode_EmptyComponentsSlice_ResultIsNil(t *testing.T) {
	// arrange
	components := []string{}

	// act
	result := pathToNode(components, nil)

	// assert
	if result != nil {
		t.Errorf("An empty component list should not be converted into a node, but the result was %q", result)
	}
}

func Test_pathToNode_NonEmptyComponentsSlice_ResultNotNil(t *testing.T) {
	// arrange
	components := []string{"root", "level 1", "level 2"}

	// act
	result := pathToNode(components, nil)

	// assert
	if result == nil {
		t.Errorf("The components %s should be converted to a node but the result was %q.", components, result)
	}
}

func Test_isValidPathComponent_PlainText_ResultIsTrue(t *testing.T) {
	// arrange
	component := "Sample Component-Näme 38712903 & 3223 dsad_123"

	// act
	result, _ := isValidPathComponent(component)

	// assert
	expected := true
	if result != expected {
		t.Errorf("The component %q is a valid component but but the function returned %v (Expected: %v).", component, result, expected)
	}
}

func Test_isValidPathComponent_EmptyString_ResultIsFalse(t *testing.T) {
	// arrange
	component := ""

	// act
	result, _ := isValidPathComponent(component)

	// assert
	expected := false
	if result != expected {
		t.Errorf("An empty string is not a valid path component but the function returned %v (Expected: %v).", result, expected)
	}
}

func Test_isValidPathComponent_ForwardSlash_ResultIsFalse(t *testing.T) {
	// arrange
	component := "a / test"

	// act
	result, _ := isValidPathComponent(component)

	// assert
	expected := false
	if result != expected {
		t.Errorf("Components with forward slashes are not valid path components but the function returned %v (Expected: %v).", result, expected)
	}
}

func Test_isValidPathComponent_Backslash_ResultIsFalse(t *testing.T) {
	// arrange
	component := `a \ test`

	// act
	result, _ := isValidPathComponent(component)

	// assert
	expected := false
	if result != expected {
		t.Errorf("Components with backslashes are not valid path components but the function returned %v (Expected: %v).", result, expected)
	}
}
