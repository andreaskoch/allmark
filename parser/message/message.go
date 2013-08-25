// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser/pattern"
	"github.com/andreaskoch/allmark/repository"
	"strings"
	"time"
)

func Parse(item *repository.Item, lines []string) (sucess bool, err error) {

	// raw markdown content
	item.RawContent = strings.TrimSpace(strings.Join(lines, "\n"))

	return true, nil
}

func getDescription(item *repository.Item) string {
	return item.MetaData.Date.Format(time.RFC850)
}

func getTitle(item *repository.Item) string {
	text := pattern.HtmlTagPattern.ReplaceAllString(item.RawContent, "")
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
