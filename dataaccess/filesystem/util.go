// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/util/fsutil"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	ReservedDirectoryNames = []string{config.FilesDirectoryName, config.MetaDataFolderName}
)

// Check if the specified directory contains an item within the range of the given max depth.
func directoryContainsItems(directory string, maxdepth int) bool {

	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		childDirectory := filepath.Join(directory, entry.Name())

		if entry.IsDir() {
			if isReservedDirectory(childDirectory) {
				continue
			}

			if maxdepth > 0 {

				// recurse
				if directoryContainsItems(childDirectory, maxdepth-1) {
					return true
				}
			}

			continue
		}

		if isMarkdownFile(childDirectory) {
			return true
		}

		continue
	}

	return false
}

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
