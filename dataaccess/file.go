// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	path   string
	parent *File
	childs []*File
}

// Creates a new root File object that has no parent.
func NewRootFile(path string, childs []*File) (*File, error) {
	return newFile(path, nil, childs)
}

// Creates a new File object that is associated with a parent File.
func NewFile(path string, parent *File, childs []*File) (*File, error) {
	return newFile(path, parent, childs)
}

func newFile(path string, parent *File, childs []*File) (*File, error) {

	normalizedPath, err := NormalizePath(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create file. Error: %s", err)
	}

	return &File{
		path:   normalizedPath,
		parent: parent,
		childs: childs,
	}, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.path)
}

func (file *File) Path() string {
	return file.path
}

func (file *File) Parent() *File {
	return file.parent
}

func (file *File) Childs() []*File {
	return file.childs
}

func (file *File) Walk(callback func(file *File)) {
	callback(file)

	for _, child := range file.childs {
		child.Walk(callback)
	}
}
