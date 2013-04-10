// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	p "github.com/andreaskoch/allmark/path"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewFileIndex(directory string) *FileIndex {

	files := getFiles(directory)

	return &FileIndex{
		Items: func() []*File {
			return files
		},

		path: directory,
	}
}

type FileIndex struct {
	Items func() []*File
	path  string
}

func (fileIndex *FileIndex) String() string {
	return fmt.Sprintf("File Index %s", fileIndex.path)
}

func (fileIndex *FileIndex) Path() string {
	return fileIndex.path
}

func (fileIndex *FileIndex) GetFilesByPath(path string, condition func(pather p.Pather) bool) []*File {

	// normalize path
	path = strings.Replace(path, p.UrlDirectorySeperator, p.FilesystemDirectorySeperator, -1)
	path = strings.Trim(path, p.FilesystemDirectorySeperator)

	// make path relative
	if strings.Index(path, FilesDirectoryName) == 0 {
		path = path[len(FilesDirectoryName):]
	}

	matchingFiles := make([]*File, 0)

	for _, file := range fileIndex.Items() {

		filePath := file.Path()
		indexPath := fileIndex.Path()

		if strings.Index(filePath, indexPath) != 0 {
			continue
		}

		relativeFilePath := filePath[len(indexPath):]
		fileMatchesPath := strings.HasPrefix(relativeFilePath, path)
		if fileMatchesPath && condition(file) {
			matchingFiles = append(matchingFiles, file)
		}
	}

	return matchingFiles
}

func getFiles(directory string) []*File {

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("Cannot read files from directory %q. Error: %s", directory, err)
		return make([]*File, 0)
	}

	files := make([]*File, 0, 5)
	for _, directoryEntry := range filesDirectoryEntries {

		// recurse
		if directoryEntry.IsDir() {
			subDirectory := filepath.Join(directory, directoryEntry.Name())
			files = append(files, getFiles(subDirectory)...)
			continue
		}

		// append new file
		filePath := filepath.Join(directory, directoryEntry.Name())
		files = append(files, NewFile(filePath))
	}

	return files
}
