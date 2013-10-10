// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"regexp"
	"strings"
)

var (
	// [*description text*](*folder path*)
	markdownLinkPattern = regexp.MustCompile(`\[(.*)\]\(([^)]+)\)`)

	markdownItemLinkPattern = regexp.MustCompile(`\[(.*)\]\(/([^)]+)\)`)
)

func Convert(item *repository.Item, pathProvider *path.Provider, itemResolver repository.ItemResolver) string {

	// files
	files := item.Files

	// assign the raw markdown content for the add-ins to work on
	convertedContent := item.RawContent

	// render markdown extensions
	convertedContent = renderImageGalleries(files, pathProvider, convertedContent)
	convertedContent = renderFileLinks(files, pathProvider, convertedContent)
	convertedContent = renderCSVTables(files, pathProvider, convertedContent)
	convertedContent = renderPDFs(files, pathProvider, convertedContent)
	convertedContent = renderVideos(files, pathProvider, convertedContent)
	convertedContent = renderAudio(files, pathProvider, convertedContent)
	convertedContent = renderFilePreview(files, pathProvider, convertedContent)

	// rewrite all links
	convertedContent = rewireLinks(files, pathProvider, convertedContent)

	// rewire item links
	convertedContent = rewireItemLinks(convertedContent, itemResolver)

	// render markdown
	convertedContent = renderMarkdown(files, pathProvider, convertedContent)

	switch itemType := item.MetaData.ItemType; itemType {
	case types.PresentationItemType:
		convertedContent = renderPresentation(convertedContent)
	}

	return convertedContent
}

func rewireItemLinks(markdown string, itemResolver repository.ItemResolver) string {

	allMatches := markdownItemLinkPattern.FindAllStringSubmatch(markdown, -1)
	for _, matches := range allMatches {

		if len(matches) != 3 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		descriptionText := strings.TrimSpace(matches[1])
		itemAlias := strings.TrimSpace(matches[2])

		// search for an item with the supplied alias
		item := itemResolver(itemAlias, anyItem)
		if item == nil {
			continue // no item found
		}

		// link
		newLinkText := fmt.Sprintf("[%s](%s)", descriptionText, item.AbsolutePath)
		markdown = strings.Replace(markdown, originalText, newLinkText, 1)

	}

	return markdown

}

func anyItem(item *repository.Item) bool {
	return true
}

func rewireLinks(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	allMatches := markdownLinkPattern.FindAllStringSubmatch(markdown, -1)
	for _, matches := range allMatches {

		if len(matches) != 3 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		descriptionText := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get matching file
		files := fileIndex.FilesByPath(path, allFiles)

		// skip if no matching files are found
		if len(files) == 0 {
			continue
		}

		// take the first file
		file := files[0]

		// link
		newLinkText := fmt.Sprintf("[%s](%s)", descriptionText, pathProvider.GetWebRoute(file))
		markdown = strings.Replace(markdown, originalText, newLinkText, 1)

	}

	return markdown
}
