// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"html"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// filepreview: [*description text*](*file path*)
	filePreviewPattern = regexp.MustCompile(`filepreview: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderFilePreview(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	for {

		found, matches := util.IsMatch(markdown, filePreviewPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the file path
		files := fileIndex.FilesByPath(path, allFiles)

		if len(files) == 0 {

			// file not found remove entry
			msg := fmt.Sprintf("<!-- Cannot embed filepreview. The file %q was not found -->", path)
			markdown = strings.Replace(markdown, originalText, msg, 1)
			continue

		}

		matchedFile := files[0]
		realFilePath := matchedFile.Path()
		fileRoute := pathProvider.GetWebRoute(matchedFile)
		displayPath := strings.TrimPrefix(path, "files/")
		linkTitle := fmt.Sprintf("%s (%s)", title, displayPath)

		// read the file
		content, contentType, err := getFileContent(realFilePath)
		if err != nil {

			// file not found remove entry
			msg := fmt.Sprintf("<!-- Cannot read file %q (Error: %s) -->", path, err)
			markdown = strings.Replace(markdown, originalText, msg, 1)
			continue

		}

		// assemble the file preview code
		previewCode := fmt.Sprintf(`<section class="filepreview filepreview-%s">`, contentType)
		previewCode += fmt.Sprintf("\n<h1><a href=\"%s\" title=\"%s\">%s</a></h1>\n", fileRoute, linkTitle, title)

		previewCode += "<pre>\n"
		if contentType != "" {
			previewCode += fmt.Sprintf(`<code class="%s">%s</code>`, contentType, content)
		} else {
			previewCode += fmt.Sprintf("<code>%s</code>", content)
		}
		previewCode += "\n</pre>"

		previewCode += fmt.Sprintf("</section>\n\n")

		// replace markdown with image gallery
		markdown = strings.Replace(markdown, originalText, previewCode, 1)

	}

	return markdown
}

func getFileContent(path string) (fileContent string, contentType string, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	// convert to string
	fileContent = string(content)

	// escape html charachters
	fileContent = html.EscapeString(fileContent)

	// get the content type
	contentType = getFileType(path)

	return fileContent, contentType, nil
}

func getFileType(path string) string {

	// get the file extension
	extension := strings.ToLower(strings.TrimSpace(filepath.Ext(path)))

	// remove the leading dot
	extension = strings.TrimPrefix(extension, ".")

	return extension
}
