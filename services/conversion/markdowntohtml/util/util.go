// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
)

func GetFallbackLink(title, path string) string {
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func GetFileContent(file *model.File) (fileContent string, contentType string, err error) {

	// get the file content
	contentProvider := file.ContentProvider()
	data, err := contentProvider.Data()
	if err != nil {
		return
	}

	fileContent = string(data) // convert data to string

	// get the mime type
	contentType, err = contentProvider.MimeType()
	if err != nil {
		return
	}

	return fileContent, contentType, nil
}
