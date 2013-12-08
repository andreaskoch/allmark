// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"testing"
)

func Test_getFallbackAlias_RouteWithFile_ReturnsItemFolder(t *testing.T) {
	// arrange
	expected := "test"
	inputPath := fmt.Sprintf("/repository/document/%s/document.md", expected)
	route, _ := route.New("/repository", inputPath)

	item, _ := model.NewItem(route, []*model.File{})

	// act
	result := getFallbackAlias(item)

	// assert
	if result != expected {
		t.Errorf("The result was expected to be %q but was %q.", expected, result)
	}
}
