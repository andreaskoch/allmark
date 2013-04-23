// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"strings"
)

func Convert(item *repository.Item, rawLines []string) string {

	// create context
	fileIndex := item.Files
	repositoryPathProvider := path.NewProvider(item.Directory())
	rawMarkdownContent := strings.TrimSpace(strings.Join(rawLines, "\n"))

	// image gallery
	galleryRenderer := NewImageGalleryRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = galleryRenderer(rawMarkdownContent)

	// tables
	tableRenderer := NewTableRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = tableRenderer(rawMarkdownContent)

	// markdown to html
	html := markdownToHtml(rawMarkdownContent)

	return html
}
