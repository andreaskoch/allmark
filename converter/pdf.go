// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

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
	// !pdf[*description text*](*some pdf file*)
	pdfPattern = regexp.MustCompile(`!pdf\[([^\]]+)\]\(([^)]+)\)`)
)

func newPDFRenderer(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func(text string) string {
	return func(text string) string {
		return renderPDF(markdown, fileIndex, pathProvider)
	}
}

func renderPDF(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

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
		fileTitle := fmt.Sprintf("%s - %s", title, getFileTitle(pdfFile))

		fileLinksCode := fmt.Sprintf(`<section class="pdf" title="%s" data-pdf="%s">
				<h1>%s</h1>

				<nav>
					<button class="prev">Previous</button>
    				<button class="next">Next</button>

    				<div class="pager">
    					Page: <span class="pager-current-page"></span> / <span class="pager-total-pagecount"></span>
    				</div>
				</nav>

				<div class="pdf-viewer">
					<canvas></canvas>
				</div>
			</section>`, fileTitle, filePath, title)

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
