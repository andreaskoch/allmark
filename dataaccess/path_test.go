// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"testing"
)

// Testing the NormalizePath function: Remove leading and trailing white space
func Test_NormalizePath_TrimEnd(t *testing.T) {
	// arrange
	inputPath := " documents/Test "
	expectedResult := "documents/Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed leading and trailing white space (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the NormalizePath function: Remove a trailing slash
func Test_NormalizePath_RemoveTrailingSlashes(t *testing.T) {
	// arrange
	inputPath := "documents/Test/"
	expectedResult := "documents/Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the trailing slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the NormalizePath function: Remove a trailing slash
func Test_NormalizePath_RemoveLeadingSlashes(t *testing.T) {
	// arrange
	inputPath := "/documents/Test/"
	expectedResult := "documents/Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have removed the leading slash (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the NormalizePath function: Replace all backslashes with forward slashes
func Test_NormalizePath_NormalizeSlashes(t *testing.T) {
	// arrange
	inputPath := "documents\\Test"
	expectedResult := "documents/Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all backslashes with forward slashes (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the NormalizePath function: Replace white space with url safe characters
func Test_NormalizePath_ReplaceWhitespaceWithUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "documents/A Test"
	expectedResult := "documents/A+Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all white space characters with url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}

// Testing the NormalizePath function: Replace all double white spaces with a single url safe character
func Test_NormalizePath_ReplaceDoubleWhitespaceWithASingleUrlSafeCharacters(t *testing.T) {
	// arrange
	inputPath := "my    documents/A  Test"
	expectedResult := "my+documents/A+Test"

	// act
	result := NormalizePath(inputPath)

	// assert
	if result != expectedResult {
		t.Errorf("Should have replaced all double white spaces with a single url safe characters (Expected: %s, Actual: %s)", expectedResult, result)
	}
}
