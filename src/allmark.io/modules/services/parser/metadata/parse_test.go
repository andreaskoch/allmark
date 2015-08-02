// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"testing"

	"allmark.io/modules/model"
)

func Test_normalizeAlias(t *testing.T) {
	// arrange
	input := "Some Alias "
	expected := "some-alias"

	// act
	result := normalizeAlias(input)

	// assert
	if result != expected {
		t.Errorf("normalizeAlias(%q) should return %q but returned %q instead", input, expected, result)
	}
}

//func parseTags(metaData *model.MetaData, lines []string) (remainingLines []string) {
func Test_parseTags(t *testing.T) {
	// arrange
	metaData := model.NewMetaData()
	lines := []string{
		"tags: tag1, tag2, tag3",
	}

	// act
	parseTags(metaData, lines)

	// assert
	if len(metaData.Tags) != 3 {
		t.Errorf("The parser should have found 3 tags but contained only %v.", len(metaData.Tags))
	}
}
