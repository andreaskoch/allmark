// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/view"
	"regexp"
	"time"
)

// Pattern which matches all HTML/XML tags
var HtmlTagPattern = regexp.MustCompile(`\<[^\>]*\>`)

func createMessageMapperFunc(parsedItem *parser.ParsedItem, relativPath, absolutePath string) *view.Model {

	return &view.Model{
		RelativeRoute: relativPath,
		AbsoluteRoute: absolutePath,
		Title:         getTitle(parsedItem),
		Description:   getDescription(parsedItem),
		Content:       parsedItem.ConvertedContent,
		LanguageTag:   getTwoLetterLanguageCode(parsedItem.MetaData.Language),
		Date:          formatDate(parsedItem.MetaData.Date),
		Type:          parsedItem.MetaData.ItemType,
	}

}

func getDescription(parsedItem *parser.ParsedItem) string {
	return parsedItem.MetaData.Date.Format(time.RFC850)
}

func getTitle(parsedItem *parser.ParsedItem) string {
	text := HtmlTagPattern.ReplaceAllString(parsedItem.ConvertedContent, "")
	excerpt := getTextExcerpt(text, 30)
	time := parsedItem.MetaData.Date.Format(time.RFC850)

	return fmt.Sprintf("%s: %s", time, excerpt)
}

func getTextExcerpt(text string, length int) string {

	if len(text) <= length {
		return text
	}

	return text[0:length] + " ..."
}
