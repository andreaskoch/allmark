// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"fmt"
	"testing"

	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
	"allmark.io/modules/model"
)

func Test_getFallbackAlias_RouteWithFile_ReturnsItemFolder(t *testing.T) {
	// arrange
	expected := "test"
	inputPath := fmt.Sprintf("/repository/document/%s", expected)
	route := route.NewFromFilePath("/repository", inputPath)

	item := model.NewItem(route, []*model.File{}, dataaccess.TypePhysical)

	// act
	result := getFallbackAlias(item)

	// assert
	if result != expected {
		t.Errorf("The result was expected to be %q but was %q.", expected, result)
	}
}

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
