// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

import(
	"github.com/andreaskoch/allmark/path"
	"path/filepath"
	"strings"
)

func getFileTitle(pather path.Pather) string {
	fileName := filepath.Base(pather.Path())
	fileExtension := filepath.Ext(pather.Path())

	// remove the file extension from the file name
	filenameWithoutExtension := fileName[0:(strings.LastIndex(fileName, fileExtension))]

	return filenameWithoutExtension
}

func isImageFile(pather path.Pather) bool {
	fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
	switch fileExtension {
	case ".png", ".gif", ".jpeg", ".jpg", ".svg", ".tiff":
		return true
	default:
		return false
	}

	panic("Unreachable")
}

func allFiles(pather path.Pather) bool {
	return true
}