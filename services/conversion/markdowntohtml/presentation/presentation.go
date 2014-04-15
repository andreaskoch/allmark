// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package presentation

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"strings"
)

func New(pathProvider paths.Pather, files []*model.File) *PresentationExtension {
	return &PresentationExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type PresentationExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *PresentationExtension) Convert(html string) (convertedContent string, conversionError error) {
	slides := strings.Split(html, "<hr />")
	presentationCode := fmt.Sprintf(`<section class="slide">%s</section>`, strings.Join(slides, `</section><section class="slide">`))
	return presentationCode, nil
}
