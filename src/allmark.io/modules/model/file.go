// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"allmark.io/modules/dataaccess"
	"strings"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	dataaccess.File
}

// IsImageFile returns true if the supplied file model is an image.
func IsImageFile(file *File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return IsImage(mimetype)
}

// IsImage returns true if the supplied mime type is an image.
func IsImage(mimetype string) bool {
	return strings.HasPrefix(mimetype, "image/")
}

// IsTextFile returns true if the supplied file model is a text file.
func IsTextFile(file *File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	if strings.HasPrefix(mimetype, "text/") {
		return true
	}

	if strings.Contains(mimetype, "json") {
		return true
	}

	if strings.Contains(mimetype, "javascript") {
		return true
	}

	if strings.Contains(mimetype, "xml") {
		return true
	}

	if strings.Contains(mimetype, "cert") {
		return true
	}

	return false
}

// IsAudioFile returns true if the supplied file model is an audio file.
func IsAudioFile(file *File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "audio/")
}

// IsVideoFile returns true if the supplied file model is a video file.
func IsVideoFile(file *File) bool {
	mimetype, err := GetMimeType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(mimetype, "video/")
}

// GetMimeType returns the mime type if the given file model.
func GetMimeType(file *File) (string, error) {
	mimetype, err := file.MimeType()
	if err != nil {
		return "", err
	}

	return mimetype, nil
}
