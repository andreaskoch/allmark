// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// pdf: [*description text*](*some pdf file*)
	pdfPattern = regexp.MustCompile(`pdf: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderPDFs(item *repository.Item, rawContent string) string {
	return convertPdfMarkdownExtension(rawContent, item.Files, item.FilePathProvider())
}

func convertPdfMarkdownExtension(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

	for {

		found, matches := util.IsMatch(markdown, pdfPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// create link list code
		files := fileIndex.FilesByPath(path, isPDFFile)
		if len(files) != 1 {
			// remove the invalid code
			markdown = strings.Replace(markdown, originalText, fmt.Sprintf("<!-- pdf-preview: could not find file %s -->", path), 1)
			continue
		}

		// create the pdf viewer code
		pdfFile := files[0]
		filePath := pathProvider.GetWebRoute(pdfFile)
		fileTitle := fmt.Sprintf("PDF: %s - %s", title, getFileTitle(pdfFile))

		fileLinksCode := fmt.Sprintf(`<section class="pdf">
				<h1>%s</h1>
				<a href="%s" target="_blank" title="%s">%s</a>
			</section>`, title, filePath, fileTitle, filePath)

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, fileLinksCode, 1)

	}

	return markdown
}

func isPDFFile(pather path.Pather) bool {
	fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
	switch fileExtension {
	case ".pdf":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
