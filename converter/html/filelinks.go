// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
)

var (
	// files: [*description text*](*folder path*)
	fileLinksPattern = regexp.MustCompile(`files: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderFileLinks(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	for {

		found, matches := util.IsMatch(markdown, fileLinksPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// create link list code
		files := fileIndex.FilesByPath(path, allFiles)

		fileLinks := getFileLinks(title, files, pathProvider)
		fileLinksCode := fmt.Sprintf(`<section class="filelinks">
				<h1>%s</h1>
				<ol>
					<li>
					%s
					</li>
				</ol>
			</section>`, title, strings.Join(fileLinks, "\n</li>\n<li>\n"))

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, fileLinksCode, 1)

	}

	return markdown
}

func getFileLinks(title string, files []*repository.File, pathProvider *path.Provider) []string {

	numberOfFiles := len(files)
	fileLinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {

		filePath := pathProvider.GetWebRoute(file)
		fileTitle := fmt.Sprintf("%s - %s (File %v of %v)", title, getFileTitle(file), index+1, numberOfFiles)
		linkText := getLinkTextFromFilePath(filePath)

		fileLinks[index] = fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, filePath, fileTitle, linkText)
	}

	return fileLinks
}
