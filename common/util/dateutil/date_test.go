// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dateutil

import (
	"testing"
	"time"
)

func Test_ParseIso8601Date_ValidIso8601Date_CorrectResultIsReturned(t *testing.T) {

	// Arrange
	var fallback time.Time
	dateString := "2013-07-26"
	expectedResult, err := time.Parse("2006-Jan-02", "2013-Jul-26")
	if err != nil {
		panic(err)
	}

	// Act
	result, err := ParseIso8601Date(dateString, fallback)

	// Assert
	if err != nil {
		t.Fail()
		t.Logf("Parsing the value '%v' returned an error even though no error was expected.", dateString)
	}

	if !result.Equal(expectedResult) {
		t.Fail()
		t.Logf("Parsing the value %q did not return the expected result %q.", result, expectedResult)
	}

}

func Test_ParseIso8601Date_ValidIso8601Dates_NoErrorIsReturned(t *testing.T) {

	// Arrange
	var fallback time.Time

	dateStrings := []string{
		"2013-02-08",
		"2013-01-01",
		"2013-12-31",
		"0001-01-01",
		"0001-12-31",
		"9999-01-01",
		"9999-12-31",
	}

	// Act
	for _, dateString := range dateStrings {
		_, err := ParseIso8601Date(dateString, fallback)

		// Assert
		if err != nil {
			t.Fail()
			t.Logf("Parsing the value '%v' returned an error even though no error was expected.", dateString)
		}

	}
}

func Test_ParseIso8601Date_ValidIso8601Date_WithValidTime_NoErrorIsReturned(t *testing.T) {

	// Arrange
	var fallback time.Time
	dateString := "2013-02-08 21:13"

	// Act
	_, err := ParseIso8601Date(dateString, fallback)

	// Assert
	if err != nil {
		t.Fail()
		t.Logf("Parsing the value '%v' returned an error even though no error was expected.", dateString)
	}
}

func Test_ParseIso8601Date_InvalidIso8601Dates_ErrorIsReturned(t *testing.T) {

	// Arrange
	var fallback time.Time
	dateStrings := []string{
		"99-02-08",
		"1-1-1",
		"2013-1-1",
		"2013-01-1",
		"2013-1-01",
		"13-01-01",
		"83-12-31",
		"21400-12-31",
	}

	// Act
	for _, dateString := range dateStrings {
		_, err := ParseIso8601Date(dateString, fallback)

		// Assert
		if err == nil {
			t.Fail()
			t.Logf("Parsing the value '%v' returned should return an error because it is not a valid date.", dateString)
		}

	}
}
