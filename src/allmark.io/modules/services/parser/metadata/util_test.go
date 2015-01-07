// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"testing"
)

func Test_GetMetaDataLines_SingleLine(t *testing.T) {
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
	result := GetMetaDataLines(inputLines)

	// assert
	if len(result) != 1 {
		t.Errorf("The resulting line slice should contain one line but contained %d lines.", len(result))
	}
}

func Test_GetMetaDataLines_SingleLineWithWhitespace(t *testing.T) {
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
	result := GetMetaDataLines(inputLines)

	// assert
	if len(result) != 3 {
		t.Errorf("The resulting line slice should contain three lines but contained %d lines.", len(result))
	}
}

func Test_GetMetaDataLines_EmptyMetaDataSection(t *testing.T) {
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
	result := GetMetaDataLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}

func Test_GetMetaDataLines_NoMetaData(t *testing.T) {
	// arrange
	inputLines := []string{
		"# Headline",
		"Description",
		"",
		"yada yada",
		"",
	}

	// act
	result := GetMetaDataLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}

func Test_GetMetaDataLines_NoInput(t *testing.T) {
	// arrange
	inputLines := []string{}

	// act
	result := GetMetaDataLines(inputLines)

	// assert
	if len(result) != 0 {
		t.Errorf("The resulting line slice should be empty but contained %d lines.", len(result))
	}
}
