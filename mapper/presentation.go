// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/view"
	"strings"
)

func createPresentationMapperFunc(parsedItem *parser.ParsedItem, relativPath, absolutePath string) *view.Model {

	return &view.Model{
		Level:         parsedItem.Level,
		RelativeRoute: relativPath,
		AbsoluteRoute: absolutePath,
		Title:         parsedItem.Title,
		Description:   parsedItem.Description,
		Content:       getPresentationContent(parsedItem.ConvertedContent),
		LanguageTag:   getTwoLetterLanguageCode(parsedItem.MetaData.Language),
		Date:          formatDate(parsedItem.MetaData.Date),
		Type:          parsedItem.MetaData.ItemType,
	}

}

func getPresentationContent(html string) string {
	slides := strings.Split(html, "<hr />")
	presentationCode := fmt.Sprintf(`<section class="slide">%s</section>`, strings.Join(slides, `</section><section class="slide">`))
	return presentationCode
}
