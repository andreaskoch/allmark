// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	ReservedDirectoryNames = []string{config.FilesDirectoryName, config.MetaDataFolderName}
)

func isReservedDirectory(path string) bool {

	if isFile, _ := fsutil.IsFile(path); isFile {
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
		if isMarkdown := isMarkdownFile(absoluteFilePath); isMarkdown {
			return true, absoluteFilePath
		}
	}

	return false, ""
}

func getChildDirectories(directory string) []string {

	directories := make([]string, 0)
	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		if !entry.IsDir() {
			continue // skip files
		}

		childDirectory := filepath.Join(directory, entry.Name())
		if isReservedDirectory(childDirectory) {
			continue // skip reserved directories
		}

		// append directory
		directories = append(directories, childDirectory)
	}

	return directories
}

func isMarkdownFile(fileNameOrPath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(fileNameOrPath))
	switch fileExtension {
	case ".md", ".markdown", ".mdown":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
