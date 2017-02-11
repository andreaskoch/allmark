// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/common/util/fsutil"
	"github.com/andreaskoch/allmark/dataaccess"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func newFileProvider(logger logger.Logger, repositoryPath string) (*fileProvider, error) {

	// abort if repoistory path does not exist
	if !fsutil.PathExists(repositoryPath) {
		return nil, fmt.Errorf("The repository path %q does not exist.", repositoryPath)
	}

	// abort if the supplied repository path is not a directory
	if isDirectory, _ := fsutil.IsDirectory(repositoryPath); !isDirectory {
		return nil, fmt.Errorf("The supplied item repository path %q is not a directory.", repositoryPath)
	}

	return &fileProvider{
		logger:         logger,
		repositoryPath: repositoryPath,
	}, nil
}

type fileProvider struct {
	logger         logger.Logger
	repositoryPath string
}

func (provider *fileProvider) GetFilesFromDirectory(itemDirectory, filesDirectory string) []dataaccess.File {

	children := make([]dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(filesDirectory)
	if err != nil {
		return children
	}

	for _, directoryEntry := range filesDirectoryEntries {

		filePath := filepath.Join(filesDirectory, directoryEntry.Name())

		// recurse if the path is a directory
		if isDir, _ := fsutil.IsDirectory(filePath); isDir {
			children = append(children, provider.GetFilesFromDirectory(itemDirectory, filePath)...)
			continue
		}

		// append new file
		file, err := createFileFromFilesystem(provider.repositoryPath, itemDirectory, filePath)
		if err != nil {
			provider.logger.Error("Unable to add file %q to index. Error: %s", filePath, err)
			continue
		}

		children = append(children, file)
	}

	return children
}

func createFileFromFilesystem(repositoryPath, itemDirectory, filePath string) (dataaccess.File, error) {

	// check if the file path is a file
	if isFile, _ := fsutil.IsFile(filePath); !isFile {
		return nil, fmt.Errorf("%q is not a file.", filePath)
	}

	parentRoute := route.NewFromFilePath(repositoryPath, itemDirectory)
	route := route.NewFromFilePath(repositoryPath, filePath)
	contentProvider, contentProviderError := newFileContentProviderWithoutChecksum(filePath, route)
	if contentProviderError != nil {
		return nil, contentProviderError
	}

	// create the file
	file := &File{
		contentProvider,
		parentRoute,
		route,
	}

	return file, nil
}
