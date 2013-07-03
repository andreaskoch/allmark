// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
	"regexp"
	"time"
)

// Pattern which matches all HTML/XML tags
var HtmlTagPattern = regexp.MustCompile(`\<[^\>]*\>`)

func createMessageMapperFunc(item *repository.Item) *view.Model {

	return &view.Model{
		Level:         item.Level,
		RelativeRoute: item.RelativePath,
		AbsoluteRoute: item.AbsolutePath,
		Title:         getTitle(item),
		Description:   getDescription(item),
		Content:       item.ConvertedContent,
		LanguageTag:   getTwoLetterLanguageCode(item.MetaData.Language),
		Date:          formatDate(item.MetaData.Date),
		Type:          item.MetaData.ItemType,
	}

}

func getDescription(item *repository.Item) string {
	return item.MetaData.Date.Format(time.RFC850)
}

func getTitle(item *repository.Item) string {
	text := HtmlTagPattern.ReplaceAllString(item.ConvertedContent, "")
	excerpt := getTextExcerpt(text, 30)
	time := item.MetaData.Date.Format(time.RFC850)

	return fmt.Sprintf("%s: %s", time, excerpt)
}

func getTextExcerpt(text string, length int) string {

	if len(text) <= length {
		return text
	}

	return text[0:length] + " ..."
}
