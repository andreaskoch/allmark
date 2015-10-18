// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/model"
	"fmt"
	"regexp"
	"strings"
)

var (
	htmlLinkPattern = regexp.MustCompile(`(src|href)="([^"]+)"`)
)

func rewireLinks(itemContentPathProvider paths.Pather, files []*model.File, html string) string {

	allMatches := htmlLinkPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 3 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		linkType := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get matching file
		matchingFile := getMatchingFiles(path, files)

		// skip if no matching files are found
		if matchingFile == nil {
			continue
		}

		// assemble the new link path
		matchingFilePath := itemContentPathProvider.Path(matchingFile.Route().Value())

		// assemble the new link
		newLinkText := fmt.Sprintf("%s=\"%s\"", linkType, matchingFilePath)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}

func getMatchingFiles(path string, files []*model.File) *model.File {
	for _, file := range files {
		if file.Route().IsMatch(path) {
			return file
		}
	}

	return nil
}
