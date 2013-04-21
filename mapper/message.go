// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
	"regexp"
	"time"
)

// Pattern which matches all HTML/XML tags
var HtmlTagPattern = regexp.MustCompile(`\<[^\>]*\>`)

func createMessageMapperFunc(pathProvider *path.Provider, targetFormat string) Mapper {
	return func(item *repository.Item) view.Model {

		parsed, err := converter.Convert(item, targetFormat)
		if err != nil {
			return view.Error(fmt.Sprintf("%s", err))
		}

		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       getTitle(parsed),
			Description: getDescription(parsed),
			Content:     parsed.ConvertedContent,
			LanguageTag: getTwoLetterLanguageCode(parsed.MetaData.Language),
		}
	}
}

func getDescription(parsedResult *parser.Result) string {
	return parsedResult.MetaData.Date.Format(time.RFC850)
}

func getTitle(parsedResult *parser.Result) string {
	text := HtmlTagPattern.ReplaceAllString(parsedResult.ConvertedContent, "")
	excerpt := getTextExcerpt(text, 30)
	time := parsedResult.MetaData.Date.Format(time.RFC850)

	return fmt.Sprintf("%s: %s", time, excerpt)
}

func getTextExcerpt(text string, length int) string {

	if len(text) <= length {
		return text
	}

	return text[0:length] + " ..."
}
