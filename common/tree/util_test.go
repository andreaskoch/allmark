// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"github.com/elWyatt/allmark/common/route"
	"testing"
)

func Test_RouteToPath(t *testing.T) {
	// arrange
	inputRoute := route.NewFromRequest("document/sample-doc/child-1")

	// act
	result := RouteToPath(inputRoute).String()

	// assert
	expected := "document > sample-doc > child-1"
	if result != expected {
		t.Errorf("The path should be %q but was %q instead.", expected, result)
	}
}
