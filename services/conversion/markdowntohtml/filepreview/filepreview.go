// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepreview

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/util"
	"html"
	"regexp"
	"strings"
)

var (
	// filepreview: [*description text*](*file path*)
	markdownPattern = regexp.MustCompile(`filepreview: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, files []*model.File) *FilePreviewExtension {
	return &FilePreviewExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type FilePreviewExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *FilePreviewExtension) Convert(markdown string) (convertedContent string, conversionError error) {

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
		renderedCode := converter.getPreviewCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)
	}

	return convertedContent, nil
}

func (converter *FilePreviewExtension) getPreviewCode(title, path string) string {

	// get the file that matches the path
	file := converter.getMatchingFile(path)
	if file != nil {

		filepath := converter.pathProvider.Path(file.Route().Value())

		if content, contentType, err := util.GetFileContent(file); err == nil {

			// escape html entities
			escapedContent := html.EscapeString(content)

			return fmt.Sprintf(`<section class="filepreview filepreview-%s">
			<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
			<pre>
				<code class="%s">%s</code>
			</pre>
			</section>`, contentType, filepath, title, title, contentType, escapedContent)
		}

	}

	// fallback
	return util.GetFallbackLink(title, path)
}

func (converter *FilePreviewExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsTextFile(file) {
			return file
		}
	}

	return nil
}
