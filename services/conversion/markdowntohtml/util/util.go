// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"strings"
)

func GetFallbackLink(title, path string) string {
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

func IsInternalLink(link string) bool {
	return !IsExternalLink(link)
}

func IsExternalLink(link string) bool {
	lowercase := strings.TrimSpace(strings.ToLower(link))
	isHttpLink := strings.HasPrefix(lowercase, "http:")
	isHttpsLink := strings.HasPrefix(lowercase, "https:")
	return isHttpLink || isHttpsLink
}

func IsImageFile(file *model.File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "image/")
}

func IsTextFile(file *model.File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "text/")
}

func IsAudioFile(file *model.File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "audio/")
}

func IsVideoFile(file *model.File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "video/")
}

func IsPDFFile(file *model.File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return mimetype == "application/pdf"
}

func GetMimeType(file *model.File) (string, error) {
	contentProvider := file.ContentProvider()
	mimetype, err := contentProvider.MimeType()
	if err != nil {
		return "", err
	}

	return mimetype, nil
}
