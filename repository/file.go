// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"path/filepath"
)

type File struct {
	path string

	rootPathProvider     *path.Provider
	relativePathProvider *path.Provider

	AbsolutePath string
	RelativePath string
}

func newFile(rootPathProvider *path.Provider, path string) (*File, error) {

	// create a path provider
	directory := filepath.Dir(path)
	relativePathProvider := rootPathProvider.New(directory)

	// create the file
	file := &File{
		path: path,

		relativePathProvider: relativePathProvider,
		rootPathProvider:     rootPathProvider,
	}

	// assign paths
	file.RelativePath = relativePathProvider.GetWebRoute(file)
	file.AbsolutePath = rootPathProvider.GetWebRoute(file)

	return file, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.RootPathProvider().GetWebRoute(file))
}

func (file *File) Path() string {
	return file.path
}

func (file *File) PathType() string {
	return path.PatherTypeFile
}

func (file *File) Directory() string {
	return filepath.Dir(file.Path())
}

func (file *File) RootPathProvider() *path.Provider {
	return file.rootPathProvider
}
