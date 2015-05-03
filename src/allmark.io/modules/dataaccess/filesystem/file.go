// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"allmark.io/modules/common/content"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/dataaccess"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type File struct {
	*content.ContentProvider

	parentRoute route.Route
	fileRoute   route.Route
}

func NewFile(fileRoute, parentRoute route.Route, contentProvider *content.ContentProvider) (*File, error) {
	return &File{
		contentProvider,
		parentRoute,
		fileRoute,
	}, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.fileRoute.Value())
}

func (file *File) Id() string {
	hash := hashutil.FromString(file.fileRoute.Value())
	return hash
}

func (file *File) Name() string {
	return fmt.Sprintf("%s", file.fileRoute.LastComponentName())
}

func (file *File) Parent() route.Route {
	return file.parentRoute
}

func (file *File) Route() route.Route {
	return file.fileRoute
}

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

	childs := make([]dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(filesDirectory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		filePath := filepath.Join(filesDirectory, directoryEntry.Name())

		// recurse if the path is a directory
		if isDir, _ := fsutil.IsDirectory(filePath); isDir {
			childs = append(childs, provider.GetFilesFromDirectory(itemDirectory, filePath)...)
			continue
		}

		// append new file
		file, err := newFile(provider.repositoryPath, itemDirectory, filePath)
		if err != nil {
			provider.logger.Error("Unable to add file %q to index. Error: %s", filePath, err)
			continue
		}

		childs = append(childs, file)
	}

	return childs
}

func newFile(repositoryPath, itemDirectory, filePath string) (dataaccess.File, error) {

	// check if the file path is a file
	if isFile, _ := fsutil.IsFile(filePath); !isFile {
		return nil, fmt.Errorf("%q is not a file.", filePath)
	}

	// parent route
	parentRoute, err := route.NewFromFilePath(repositoryPath, itemDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a parent route for the File with the file path %q. Error: %s", filePath, err)
	}

	// route
	route, err := route.NewFromFilePath(repositoryPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a File for the file path %q. Error: %s", filePath, err)
	}

	// content provider
	contentProvider := newFileContentProviderWithoutChecksum(filePath, route)

	// create the file
	file, err := NewFile(route, parentRoute, contentProvider)

	if err != nil {
		return nil, err
	}

	return file, nil
}
