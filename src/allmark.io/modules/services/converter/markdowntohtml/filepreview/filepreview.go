// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepreview

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/util"
	"bufio"
	"bytes"
	"fmt"
	"html"
	"io"
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

func (converter *FilePreviewExtension) Convert(markdown string) (convertedContent string, converterError error) {

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

		// determine the content type
		contentType, err := file.MimeType()
		if err != nil {
			// could not determine the mime type
			return util.GetHtmlLinkCode(title, path)
		}

		// prepare reading the file data
		bytesBuffer := new(bytes.Buffer)
		dataWriter := bufio.NewWriter(bytesBuffer)
		contentReader := func(content io.ReadSeeker) error {
			_, err := io.Copy(dataWriter, content)
			return err
		}

		if err := file.Data(contentReader); err == nil {

			// escape html entities
			escapedContent := html.EscapeString(bytesBuffer.String())

			return fmt.Sprintf(`<section class="filepreview filepreview-%s">
			<header><a href="%s" target="_blank" title="%s">%s</a></header>
			<pre>
				<code class="%s">%s</code>
			</pre>
			</section>`, contentType, filepath, title, title, contentType, escapedContent)
		}

	}

	// fallback
	return util.GetHtmlLinkCode(title, path)
}

func (converter *FilePreviewExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsTextFile(file) {
			return file
		}
	}

	return nil
}
