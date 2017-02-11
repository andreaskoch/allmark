// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
)

func Test_IsChildOf_ParentRouteIsEmpty_ResultIsTrue(t *testing.T) {
	// arrange
	parent := NewFromRequest("/")
	child := NewFromRequest("/documents/Test-1")

	// act
	result := child.IsChildOf(parent)

	// assert
	if !result {
		t.Errorf("%q is a direct child of %q. The result should be true but was %t.", child, parent, result)
	}
}

func Test_IsChildOf_RouteIsFirstLevelChild_ResultIsTrue(t *testing.T) {
	// arrange
	parent := NewFromRequest("/documents/Collection")
	child := NewFromRequest("/documents/Collection/Level-1")

	// act
	result := child.IsChildOf(parent)

	// assert
	if !result {
		t.Errorf("%q is a 1st level child of %q. The result should be true but was %t.", child, parent, result)
	}
}

func Test_IsChildOf_RouteIsSecondLevelChild_ResultIsTrue(t *testing.T) {
	// arrange
	parent := NewFromRequest("/documents/Collection")
	child := NewFromRequest("/documents/Collection/Level-1/Level-2")

	// act
	result := child.IsChildOf(parent)

	// assert
	if !result {
		t.Errorf("%q is a 2nd level child of %q. The result should be true but was %t.", child, parent, result)
	}
}

func Test_IsChildOf_RouteIsNotAChild_ResultIsFalse(t *testing.T) {
	// arrange
	parent := NewFromRequest("/documents/Collection")
	child := NewFromRequest("/pictures/Test-1")

	// act
	result := child.IsChildOf(parent)

	// assert
	if result {
		t.Errorf("%q is not a child of %q. The result should be false but was %t.", child, parent, result)
	}
}
