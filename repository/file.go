// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

type File struct {
	path string
}

func NewFile(filePath string) *File {
	return &File{
		path: filePath,
	}
}

func (file *File) Path() string {
	return file.path
}
func (file *File) PathType() string {
	return "file"
}
