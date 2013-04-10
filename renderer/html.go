// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/russross/blackfriday"
	"strings"
)

func html(item *repository.Item) string {

	// render image galleries
	renderImageGalleries(item)

	// convert markdown to html
	rawMarkdownContent := strings.TrimSpace(strings.Join(item.RawLines, "\n"))
	html := markdownToHtml(rawMarkdownContent)

	return html
}

func markdownToHtml(markdown string) (html string) {
	return string(blackfriday.MarkdownCommon([]byte(markdown)))
}
