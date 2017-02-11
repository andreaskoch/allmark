// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"testing"
)

func Test_IsAbsoluteURI_ParameterIsRelativeURI_ResultIsFalse(t *testing.T) {
	// arrange
	input := "/yadayada.html"
	expected := false

	// act
	result := IsAbsoluteURI(input)

	// assert
	if result != expected {
		t.Errorf("The result for IsAbsoluteURI(%q) should be %v but was %v.", input, expected, result)
	}
}

func Test_IsAbsoluteURI_ParameterIsHTTPURL_ResultIsTrue(t *testing.T) {
	// arrange
	input := "http://example.com/yadayada.html"
	expected := true

	// act
	result := IsAbsoluteURI(input)

	// assert
	if result != expected {
		t.Errorf("The result for IsAbsoluteURI(%q) should be %v but was %v.", input, expected, result)
	}
}

func Test_IsAbsoluteURI_ParameterIsHTTPURL_Uppercase_ResultIsTrue(t *testing.T) {
	// arrange
	input := "HTTP://example.com/yadayada.html"
	expected := true

	// act
	result := IsAbsoluteURI(input)

	// assert
	if result != expected {
		t.Errorf("The result for IsAbsoluteURI(%q) should be %v but was %v.", input, expected, result)
	}
}

func Test_IsAbsoluteURI_ParameterIsHTTPsURL_ResultIsTrue(t *testing.T) {
	// arrange
	input := "https://example.com/yadayada.html"
	expected := true

	// act
	result := IsAbsoluteURI(input)

	// assert
	if result != expected {
		t.Errorf("The result for IsAbsoluteURI(%q) should be %v but was %v.", input, expected, result)
	}
}

func Test_IsAbsoluteURI_ParametersAreAbsoluteURLs_ResultIsTrue(t *testing.T) {
	// arrange
	inputs := []string{
		"ftp://example.com",
		"sftp://example.com",
		"ssh://example.com",
		"bitcoin:example.com",
		"mailto:jsmith@example.com",
	}
	expected := true

	for _, input := range inputs {
		// act
		result := IsAbsoluteURI(input)

		// assert
		if result != expected {
			t.Errorf("The result for IsAbsoluteURI(%q) should be %v but was %v.", input, expected, result)
		}
	}
}
