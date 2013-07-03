// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"strings"
)

func ToHtml(item *repository.Item, rawLines []string) *repository.Item {

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

	// csv files
	csvRenderer := newCSVRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = csvRenderer(rawMarkdownContent)

	// pdf
	pdfRenderer := newPDFRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = pdfRenderer(rawMarkdownContent)

	// video
	videoRenderer := newVideoRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = videoRenderer(rawMarkdownContent)

	// audio
	audioRenderer := newAudioRenderer(rawMarkdownContent, fileIndex, repositoryPathProvider)
	rawMarkdownContent = audioRenderer(rawMarkdownContent)

	// markdown to html
	html := markdown.ToHtml(rawMarkdownContent)

	item.ConvertedContent = html

	return item
}
