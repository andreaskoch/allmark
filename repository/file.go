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
}

func newFile(rootPathProvider *path.Provider, path string) (*File, error) {

	// create the file
	file := &File{
		path: path,

		rootPathProvider: rootPathProvider,
	}

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
