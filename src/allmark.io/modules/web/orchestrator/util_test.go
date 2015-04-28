// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"testing"
	"time"
)

func Test_getFormattedDate_DateIsZero_ReturnsEmptyString(t *testing.T) {
	// arrange
	inputDate := time.Time{}
	expected := ""

	// act
	result := getFormattedDate(inputDate)

	// assert
	if result != expected {
		t.Errorf("The result of getFormattedDate(%q) should be %q but was %q.", inputDate, expected, result)
	}
}
