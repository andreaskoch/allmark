// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filepreview

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pattern"
	"html"
	"regexp"
	"strings"
)

var (
	// filepreview: [*description text*](*file path*)
	filePreviewPattern = regexp.MustCompile(`filepreview: \[([^\]]+)\]\(([^)]+)\)`)
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

		found, matches := pattern.IsMatch(convertedContent, filePreviewPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the file that matches the path
		file := converter.getMatchingFile(path)
		if file != nil {

			filepath := converter.pathProvider.Path(file.Route().Value())

			// get the code
			renderedCode := getPreviewCode(title, filepath, file)

			// replace markdown
			convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

		} else {

			// fallback
			fallback := getFallback(title, path)
			convertedContent = strings.Replace(convertedContent, originalText, fallback, 1)

		}

	}

	return convertedContent, nil
}

func (converter *FilePreviewExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) {
			return file
		}
	}

	return nil
}

func getPreviewCode(title, path string, file *model.File) string {

	if content, contentType, err := getFileContent(file); err == nil {
		return fmt.Sprintf(`<section class="filepreview filepreview-%s">
			<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
			<pre>
				<code class="%s">%s</code>
			</pre>
		</section>`, contentType, path, title, title, contentType, content)
	}

	// fallback
	return getFallback(title, path)
}

func getFallback(title, path string) string {
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func getFileContent(file *model.File) (fileContent string, contentType string, err error) {

	// get the file content
	contentProvider := file.ContentProvider()
	data, err := contentProvider.Data()
	if err != nil {
		return
	}

	fileContent = string(data)                   // convert data to string
	fileContent = html.EscapeString(fileContent) // escape html charachters

	// get the mime type
	contentType, err = contentProvider.MimeType()
	if err != nil {
		return
	}

	return fileContent, contentType, nil
}
