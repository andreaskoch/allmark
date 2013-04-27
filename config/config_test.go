// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"testing"
)

func Test_getUserHomeDir_ParsedItemIsNotEmpty_ErrorIsNil(t *testing.T) {
	// act
	result, err := getUserHomeDir()

	// assert: result not empty
	if len(result) == 0 {
		t.Fail()
		t.Logf("The result should not be empty")
	}

	// assert: error is nil
	if err != nil {
		t.Fail()
		t.Logf("The function should not return an error but returned %s", err)
	}
}
