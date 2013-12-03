// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"testing"
)

func Test_GetLines_SingleLine(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",
		"Description",
		"",
		"yada yada",
		"",
		"---",
		"type: document",
	}

	// act
	result := GetLines(inputLines)

	// assert
	if len(result) != 1 {
		t.Errorf("The resulting line slice should contain one line but contained %d lines.", len(result))
	}
}

func Test_GetLines_SingleLineWithWhitespace(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",
		"Description",
		"",
		"yada yada",
		"",
		"---",
		"",
		"type: document",
		"",
	}

	// act
	result := GetLines(inputLines)

	// assert
	if len(result) != 3 {
		t.Errorf("The resulting line slice should contain three lines but contained %d lines.", len(result))
	}
}

func Test_GetLines_EmptyMetaDataSection(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",
		"Description",
		"",
		"yada yada",
		"",
		"---",
	}

	// act
	result := GetLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}

func Test_GetLines_NoMetaData(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",
		"Description",
		"",
		"yada yada",
		"",
	}

	// act
	result := GetLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}

func Test_GetLines_NoInput(t *testing.T) {
	// arrange
	inputLines := []string{}

	// act
	result := GetLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}

func Test_GetLocation_SingleLine(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",     // Line number: 0
		"Description",    // Line number: 1
		"",               // Line number: 2
		"yada yada",      // Line number: 3
		"",               // Line number: 4
		"---",            // Line number: 5
		"type: document", // Line number: 6
	}
	expectedResult := 5

	// act
	result, _ := GetLocation(inputLines)

	// assert
	if result != expectedResult {
		t.Errorf("The location of the meta data should be %d but was %d.", expectedResult, result)
	}
}
