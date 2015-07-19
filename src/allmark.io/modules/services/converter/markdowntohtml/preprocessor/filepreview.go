// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/util"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	// filepreview: [*description text*](*file path*)
	filePreviewMarkdownExtension = regexp.MustCompile(`filepreview: \[([^\]]+)\]\(([^)]+)\)`)
)

func newFilePreviewExtension(pathProvider paths.Pather, files []*model.File) *filePreviewExtension {
	return &filePreviewExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type filePreviewExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *filePreviewExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, filePreviewMarkdownExtension)
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

func (converter *filePreviewExtension) getPreviewCode(title, path string) string {

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

			code := fmt.Sprintf("**[%s](%s)**\n\n", title, filepath)
			code += fmt.Sprintf("```%s\n", contentType)
			code += strings.TrimSpace(bytesBuffer.String()) + "\n"
			code += "```"

			return code
		}

	}

	// fallback
	return util.GetHtmlLinkCode(title, path)
}

func (converter *filePreviewExtension) getMatchingFile(path string) *model.File {

	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsTextFile(file) {
			return file
		}
	}

	return nil
}
