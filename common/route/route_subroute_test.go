// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"testing"
)

func Test_SubRoute_TwoLevelRoute_RequestLevelZero_LevelZeroIsReturned(t *testing.T) {
	// arrange
	level := 0
	expected := ""
	route, _ := NewFromRequest("/documents/Test-1")

	// act
	result, _ := route.SubRoute(level)

	// assert
	if result == nil || result.Value() != expected {
		t.Errorf("The level-%d sub route should be %q, but was %q.", level, expected, result.Value())
	}
}

func Test_SubRoute_TwoLevelRoute_RequestFirstLevel_FirstLevelIsReturned(t *testing.T) {
	// arrange
	level := 1
	expected := "documents"
	route, _ := NewFromRequest("/documents/Test-1")

	// act
	result, _ := route.SubRoute(level)

	// assert
	if result == nil || result.Value() != expected {
		t.Errorf("The level-%d sub route should be %q, but was %q.", level, expected, result.Value())
	}
}

func Test_SubRoute_TwoLevelRoute_RequestSecondLevel_SecondLevelIsReturned(t *testing.T) {
	// arrange
	level := 2
	expected := "documents/Test-1"
	route, _ := NewFromRequest("/documents/Test-1")

	// act
	result, _ := route.SubRoute(level)

	// assert
	if result == nil || result.Value() != expected {
		t.Errorf("The level-%d sub route should be %q, but was %q.", level, expected, result.Value())
	}
}
