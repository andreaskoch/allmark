// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"testing"
)

func Test_GetMetaDataPosition_SingleLine(t *testing.T) {
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
	result, _ := GetMetaDataPosition(inputLines)

	// assert
	if result != expectedResult {
		t.Errorf("The location of the meta data should be %d but was %d.", expectedResult, result)
	}
}
