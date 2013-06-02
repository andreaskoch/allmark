// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

import (
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"strings"
)

func toHtml(item *repository.Item, rawLines []string) string {

	// create context
	fileIndex := item.Files
	repositoryPathProvider := path.NewProvider(item.Directory(), false)
	rawMarkdownContent := strings.TrimSpace(strings.Join(rawLines, "\n"))

	// image gallery
	galleryRenderer := newImageGalleryRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = galleryRenderer(rawMarkdownContent)

	// file links
	fileLinksRenderer := newFileLinksRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = fileLinksRenderer(rawMarkdownContent)

	// tables
	tableRenderer := newTableRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = tableRenderer(rawMarkdownContent)

	// markdown to html
	html := markdown.ToHtml(rawMarkdownContent)

	return html
}
