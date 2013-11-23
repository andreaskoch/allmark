// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
)

// Testing the New Route functon: The constructor function should return an error if the supplied path is empty.
func Test_New_EmptyStringReturnsError(t *testing.T) {
	// act
	_, err := New("")

	// assert
	if err == nil {
		t.Errorf("The constructor function should return an error if the supplied path is empty.")
	}
}

// Testing the New Route functon: The constructor function should return a Route object for the supplied path.
func Test_New_ValidPathReturnsRoute(t *testing.T) {
	// arrange
	inputPath := "document/Test"

	// act
	result, err := New(inputPath)

	// assert
	if result == nil || err != nil {
		t.Errorf("The constructor function should return a Route object for the path %q.", inputPath)
	}
}

// Testing the normalizePath function: The function returns an error if the supplied path is empty.
func Test_normalizePath_EmptyStringCreatesError(t *testing.T) {
	// arrange
	inputPath := " "

	// act
	_, err := normalizePath(inputPath)

	// assert
	if err == nil {
		t.Errorf("The have returned an error. %q is not a valid path.", inputPath)
	}
}

// Testing the normalizePath function: Remove leading and trailing white space.
func Test_normalizePath_TrimEnd(t *testing.T) {
	// arrange
	inputPath := " documents/Test "
	expectedResult := "documents/Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed leading and trailing white space (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalizePath function: Remove a trailing slash.
func Test_normalizePath_RemoveTrailingSlashes(t *testing.T) {
	// arrange
	inputPath := "documents/Test/"
	expectedResult := "documents/Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the trailing slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalizePath function: Remove a trailing slash.
func Test_normalizePath_RemoveLeadingSlashes(t *testing.T) {
	// arrange
	inputPath := "/documents/Test/"
	expectedResult := "documents/Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the leading slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalizePath function: Replace all backslashes with forward slashes.
func Test_normalizePath_NormalizeSlashes(t *testing.T) {
	// arrange
	inputPath := "documents\\Test"
	expectedResult := "documents/Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all backslashes with forward slashes (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalizePath function: Replace white space with url safe characters.
func Test_normalizePath_ReplaceWhitespaceWithUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "documents/A Test"
	expectedResult := "documents/A+Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all white space characters with url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalizePath function: Replace all double white spaces with a single url safe character.
func Test_normalizePath_ReplaceDoubleWhitespaceWithASingleUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "my    documents/A  Test"
	expectedResult := "my+documents/A+Test"

	// act
	result, _ := normalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all double white spaces with a single url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}
