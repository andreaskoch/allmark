// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"testing"
)

func Test_AbsoluteWebPathProvider_NoPrefix_Path_ReturnsPathWithoutModification(t *testing.T) {
	// arrange
	prefix := ""
	pathProvider := newAbsoluteWebPathProvider(prefix)
	inputPath := "yada/yada.html"
	expected := "yada/yada.html"

	// act
	result := pathProvider.Path(inputPath)

	// assert
	if result != expected {
		t.Errorf("The result for pathProvider.Path(%q) with a prefix of %q should be %q but was %q.", inputPath, prefix, expected, result)
	}
}

func Test_AbsoluteWebPathProvider_SlashPrefix_Path_ReturnsPathWithSlash(t *testing.T) {
	// arrange
	prefix := "/"
	pathProvider := newAbsoluteWebPathProvider(prefix)
	inputPath := "yada/yada.html"
	expected := "/yada/yada.html"

	// act
	result := pathProvider.Path(inputPath)

	// assert
	if result != expected {
		t.Errorf("The result for pathProvider.Path(%q) with a prefix of %q should be %q but was %q.", inputPath, prefix, expected, result)
	}
}

func Test_AbsoluteWebPathProvider_Path_ParameterIsAbsolute_ReturnsPathWithoutModification(t *testing.T) {
	// arrange
	prefix := "/"
	pathProvider := newAbsoluteWebPathProvider(prefix)
	inputPath := "http://example.com/yada/yada.html"
	expected := "http://example.com/yada/yada.html"

	// act
	result := pathProvider.Path(inputPath)

	// assert
	if result != expected {
		t.Errorf("The result for pathProvider.Path(%q) with a prefix of %q should be %q but was %q.", inputPath, prefix, expected, result)
	}
}
