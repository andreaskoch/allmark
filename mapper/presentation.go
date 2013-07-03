// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
	"strings"
)

func createPresentationMapperFunc(item *repository.Item) *view.Model {

	return &view.Model{
		Level:         item.Level,
		RelativeRoute: item.RelativePath,
		AbsoluteRoute: item.AbsolutePath,
		Title:         item.Title,
		Description:   item.Description,
		Content:       getPresentationContent(item.ConvertedContent),
		LanguageTag:   getTwoLetterLanguageCode(item.MetaData.Language),
		Date:          formatDate(item.MetaData.Date),
		Type:          item.MetaData.ItemType,
	}

}

func getPresentationContent(html string) string {
	slides := strings.Split(html, "<hr />")
	presentationCode := fmt.Sprintf(`<section class="slide">%s</section>`, strings.Join(slides, `</section><section class="slide">`))
	return presentationCode
}
