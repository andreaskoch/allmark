// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"strings"
)

func NewConverter(item *repository.Item) func() string {

	// create context
	fileIndex := item.Files
	pathProvider := path.NewProvider(item.Directory())
	rawMarkdownContent := strings.TrimSpace(strings.Join(item.RawLines, "\n"))

	return func() string {

		markdown := rawMarkdownContent

		// image gallery
		galleryRenderer := NewImageGalleryRenderer(rawMarkdownContent, fileIndex, pathProvider)
		markdown = galleryRenderer(markdown)

		// markdown to html
		html := markdownToHtml(markdown)

		return html
	}
}
