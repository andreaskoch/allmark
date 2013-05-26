// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/util"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	ReservedDirectoryNames = []string{FilesDirectoryName, config.MetaDataFolderName}
)

func isReservedDirectory(path string) bool {

	if isFile, _ := util.IsFile(path); isFile {
		path = filepath.Dir(path)
	}

	// get the directory name
	directoryName := strings.ToLower(filepath.Base(path))

	// all dot-directories are ignored
	if strings.HasPrefix(directoryName, ".") {
		return true
	}

	// check the reserved directory names
	for _, reservedDirectoryName := range ReservedDirectoryNames {
		if directoryName == strings.ToLower(reservedDirectoryName) {
			return true
		}
	}

	return false
}

func findMarkdownFileInDirectory(directory string) (found bool, file string) {
	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		return false, ""
	}

	for _, element := range entries {

		if element.IsDir() {
			continue // skip directories
		}

		absoluteFilePath := filepath.Join(directory, element.Name())
		if isMarkdown := markdown.IsMarkdownFile(absoluteFilePath); isMarkdown {
			return true, absoluteFilePath
		}
	}

	return false, ""
}
