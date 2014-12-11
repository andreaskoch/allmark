// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashutil

import (
	"bytes"
	"testing"
)

func Test_GetHash_ResultIsNotEmpty(t *testing.T) {

	// arrange
	inputString := "La di da"
	input := bytes.NewReader([]byte(inputString))

	// act
	result, _ := GetHash(input)

	// assert
	if result == "" {
		t.Errorf("The GetHash function should be able to calculate a hash for the supplied text %q, but the result was empty.", inputString)
	}
}

func Test_GetHash_CorrectSHA1IsReturned(t *testing.T) {

	// arrange
	inputString := "La di da"
	input := bytes.NewReader([]byte(inputString))
	expectedResult := "8-18889A25"

	// act
	result, _ := GetHash(input)

	// assert
	if result != expectedResult {
		t.Errorf("The GetHash function should return the correct hash for the string %q. (Expected: %q, Actual: %q)", inputString, expectedResult, result)
	}
}
