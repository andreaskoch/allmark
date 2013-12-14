// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"testing"
)

func Test_IsParentOf_RouteIsFirstLevelChild_ResultIsTrue(t *testing.T) {
	// arrange
	parent, _ := NewFromRequest("/documents/Collection")
	child, _ := NewFromRequest("/documents/Collection/Level-1")

	// act
	result := parent.IsParentOf(child)

	// assert
	if !result {
		t.Errorf("%q is a 1st level parent of %q. The result should be true but was %t.", child, parent, result)
	}
}

func Test_IsParentOf_RouteIsSecondLevelChild_ResultIsFalse(t *testing.T) {
	// arrange
	parent, _ := NewFromRequest("/documents/Collection")
	child, _ := NewFromRequest("/documents/Collection/Level-1/Level-2")

	// act
	result := parent.IsParentOf(child)

	// assert
	if result {
		t.Errorf("%q is only a 2nd level parent of %q. The result should be false but was %t.", child, parent, result)
	}
}

func Test_IsParentOf_RouteIsNotAParent_ResultIsFalse(t *testing.T) {
	// arrange
	parent, _ := NewFromRequest("/documents/Collection")
	child, _ := NewFromRequest("/pictures/Test-1")

	// act
	result := parent.IsParentOf(child)

	// assert
	if result {
		t.Errorf("%q is not a parent of %q. The result should be false but was %t.", child, parent, result)
	}
}
