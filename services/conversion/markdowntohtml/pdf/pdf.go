// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pdf

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// pdf: [*description text*](*some pdf file*)
	pdfPattern = regexp.MustCompile(`pdf: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, files []*model.File) *PDFExtension {
	return &PDFExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type PDFExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *PDFExtension) Convert(markdown string) (convertedContent string, conversionError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, pdfPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// fix the path
		path = converter.pathProvider.Path(path)

		// get the code
		renderedCode := getPDFCode(title, path)

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func getPDFCode(title, path string) string {

	// html5 audio file
	if isPDFFileLink(path) {
		return getPDFFileLink(title, path)
	}

	// fallback
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func isPDFFileLink(link string) (isPDFFile bool) {
	extension := filepath.Ext(link)
	extension = strings.ToLower(extension)

	switch extension {
	case ".pdf":
		return true
	default:
		return false
	}

	panic("Unreachable")
}

func getPDFFileLink(title, link string) string {
	return fmt.Sprintf(`<section class="pdf">
				<h1>%s</h1>
				<a href="%s" target="_blank" title="%s">%s</a>
			</section>`, title, link, title, link)
}
