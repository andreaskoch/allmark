// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
)

func Test_Level_RootItemRoute_LevelIsZero(t *testing.T) {
	// arrange
	route := NewFromRequest("/")

	// act
	result := route.Level()

	// assert
	expected := 0
	if result != expected {
		t.Errorf("The level of %q should be %d but was %d.", route, expected, result)
	}
}

func Test_Level_FirstLevelRoute_LevelIsOne(t *testing.T) {
	// arrange
	route := NewFromRequest("/documents")

	// act
	result := route.Level()

	// assert
	expected := 1
	if result != expected {
		t.Errorf("The level of %q should be %d but was %d.", route, expected, result)
	}
}

func Test_Level_SecondLevelRoute_LevelIsTwo(t *testing.T) {
	// arrange
	route := NewFromRequest("/documents/Test-123")

	// act
	result := route.Level()

	// assert
	expected := 2
	if result != expected {
		t.Errorf("The level of %q should be %d but was %d.", route, expected, result)
	}
}

func Test_Level_ThirdLevelRoute_LevelIsThree(t *testing.T) {
	// arrange
	route := NewFromRequest("/documents/Test-123/Another-Test")

	// act
	result := route.Level()

	// assert
	expected := 3
	if result != expected {
		t.Errorf("The level of %q should be %d but was %d.", route, expected, result)
	}
}
