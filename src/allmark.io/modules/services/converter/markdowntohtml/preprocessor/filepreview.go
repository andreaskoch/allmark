// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/util"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
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

	for _, match := range filePreviewMarkdownExtension.FindAllStringSubmatch(convertedContent, -1) {

		if len(match) != 3 {
			continue
		}

		// parameters
		originalText := strings.TrimSpace(match[0])
		title := strings.TrimSpace(match[1])
		path := strings.TrimSpace(match[2])

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
		contentLanguage := getContentLanguageFromFile(file)

		// prepare reading the file data
		bytesBuffer := new(bytes.Buffer)
		dataWriter := bufio.NewWriter(bytesBuffer)
		contentReader := func(content io.ReadSeeker) error {
			_, err := io.Copy(dataWriter, content)
			return err
		}

		if err := file.Data(contentReader); err == nil {

			code := fmt.Sprintf("**[%s](%s)**\n\n", title, filepath)
			code += fmt.Sprintf("```%s\n", contentLanguage)
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
		if file.Route().IsMatch(path) && model.IsTextFile(file) {
			return file
		}
	}

	return nil
}

// getContentLanguageFromFile derives the file content language (e.g. go, php, js, ...)
func getContentLanguageFromFile(file *model.File) string {
	return filepath.Ext(file.Route().Value())
}
