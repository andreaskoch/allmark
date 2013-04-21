// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	p "github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/watcher"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func NewFileIndex(directory string) (*FileIndex, error) {

	files := getFiles(directory)
	getFilesFunc := func() []*File {
		return files
	}

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(directory)
	if err != nil {
		return nil, fmt.Errorf("Could not create a change handler for folder %q.\nError: %s\n", directory, err)
	}

	fileIndex := &FileIndex{
		ChangeHandler: changeHandler,
		Items:         getFilesFunc,

		path: directory,
	}

	fileIndex.OnChange("Reindex files on change", func(event *watcher.WatchEvent) {
		fmt.Println("Reindexing")
		files := getFiles(directory)

		fileIndex.Items = func() []*File {
			return files
		}
	})

	return fileIndex, nil
}

type FileIndex struct {
	*watcher.ChangeHandler
	Items func() []*File

	path string
}

func (fileIndex *FileIndex) String() string {
	return fmt.Sprintf("%s", fileIndex.path)
}

func (fileIndex *FileIndex) Path() string {
	return fileIndex.path
}

func (fileIndex *FileIndex) Directory() string {
	return fileIndex.Path()
}

func (fileIndex *FileIndex) PathType() string {
	return p.PatherTypeIndex
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

	files := make([]*File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return files
	}

	for _, directoryEntry := range filesDirectoryEntries {

		// recurse
		if directoryEntry.IsDir() {
			subDirectory := filepath.Join(directory, directoryEntry.Name())
			files = append(files, getFiles(subDirectory)...)
			continue
		}

		// append new file
		filePath := filepath.Join(directory, directoryEntry.Name())
		file, err := NewFile(filePath)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", filePath, err)
		}

		files = append(files, file)
	}

	return files
}
