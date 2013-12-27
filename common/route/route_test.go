// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
)

// Testing the NewFromRequest Route functon: The constructor function should return a valid route if the supplied path is empty.
func Test_NewFromRequest_EmptyString_ValidRouteIsReturned(t *testing.T) {
	// act
	result, _ := NewFromRequest("")

	// assert
	if result == nil {
		t.Errorf("The constructor function should returned a valid route but returned nil instead.")
	}
}

// Testing the NewFromRequest Route functon: The constructor function should return a Route object for the supplied path.
func Test_NewFromRequest_ValidPathReturnsRoute(t *testing.T) {
	// arrange
	inputPath := "document/Test"

	// act
	result, err := NewFromRequest(inputPath)

	// assert
	if result == nil || err != nil {
		t.Errorf("The constructor function should return a Route object for the path %q.", inputPath)
	}
}

// Testing the normalize function: The function returns an empty route.
func Test_normalize_EmptyString_RouteIsReturned(t *testing.T) {
	// arrange
	inputPath := " "
	expected := ""

	// act
	result := normalize(inputPath)

	// assert
	if result != expected {
		t.Errorf("The normalze function should have returned %q but retured %q instead.", expected, result)
	}
}

// Testing the normalize function: Remove leading and trailing white space.
func Test_normalize_TrimEnd(t *testing.T) {
	// arrange
	inputPath := " documents/Test "
	expectedResult := "documents/Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed leading and trailing white space (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalize function: Remove a trailing slash.
func Test_normalize_RemoveTrailingSlashes(t *testing.T) {
	// arrange
	inputPath := "documents/Test/"
	expectedResult := "documents/Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the trailing slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalize function: Remove a trailing slash.
func Test_normalize_RemoveLeadingSlashes(t *testing.T) {
	// arrange
	inputPath := "/documents/Test/"
	expectedResult := "documents/Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the leading slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalize function: Replace all backslashes with forward slashes.
func Test_normalize_NormalizeSlashes(t *testing.T) {
	// arrange
	inputPath := "documents\\Test"
	expectedResult := "documents/Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all backslashes with forward slashes (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalize function: Replace white space with url safe characters.
func Test_normalize_ReplaceWhitespaceWithUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "documents/A Test"
	expectedResult := "documents/A+Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all white space characters with url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the normalize function: Replace all double white spaces with a single url safe character.
func Test_normalize_ReplaceDoubleWhitespaceWithASingleUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "my    documents/A  Test"
	expectedResult := "my+documents/A+Test"

	// act
	result := normalize(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all double white spaces with a single url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}
