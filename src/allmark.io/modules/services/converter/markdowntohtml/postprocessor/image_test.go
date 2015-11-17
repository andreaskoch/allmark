// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"fmt"
	"testing"

	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/imageprovider"
	"allmark.io/modules/services/thumbnail"
)

func Test_Convert(t *testing.T) {
	// arrange
	title := "Build Status"
	imagePath := "https://travis-ci.org/andreaskoch/allmark.png"
	expected := fmt.Sprintf(`<img src="%s" alt="%s"/>`, imagePath, title)
	input := fmt.Sprintf(`<img src="%s" alt="%s"/>`, imagePath, title)

	pathProvider := DummyPather{}
	baseRoute := route.New()
	files := []*model.File{}
	thumbnailIndex := thumbnail.EmptyIndex()
	imageProvider := imageprovider.NewImageProvider(pathProvider, thumbnailIndex)

	postprocessor := newImagePostprocessor(pathProvider, baseRoute, files, imageProvider)

	// act
	result, _ := postprocessor.Convert(input)

	// assert
	if result != expected {
		t.Errorf("The result should be %q but was %s", expected, result)
	}
}

type DummyPather struct {
}

func (p DummyPather) Path(itemPath string) string {
	return itemPath
}

func (p DummyPather) Base() route.Route {
	return route.New()
}
