// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pdf

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/util"
	"regexp"
	"strings"
)

var (
	// pdf: [*description text*](*some pdf file*)
	markdownPattern = regexp.MustCompile(`pdf: \[([^\]]+)\]\(([^)]+)\)`)
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

		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := converter.getPDFCode(title, path)

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *PDFExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsPDFFile(file) {
			return file
		}
	}

	return nil
}

func (converter *PDFExtension) getPDFCode(title, path string) string {

	if pdfFile := converter.getMatchingFile(path); pdfFile != nil {

		filepath := converter.pathProvider.Path(pdfFile.Route().Value())
		return getPDFFileLink(title, filepath)

	}

	// fallback
	return util.GetFallbackLink(title, path)
}

func getPDFFileLink(title, link string) string {
	return fmt.Sprintf(`<section class="pdf">
				<header>%s</header>
				<a href="%s" target="_blank" title="%s">%s</a>
			</section>`, title, link, title, link)
}
