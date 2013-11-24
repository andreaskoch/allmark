// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
)

// A File represents a file ressource that is associated with an Item.
type File struct {
	route  *route.Route
	parent *File
	childs []*File
}

// Creates a new root File object that has no parent.
func NewRootFile(path string) (*File, error) {
	return newFile(path, nil)
}

// Creates a new File object that is associated with a parent File.
func NewFile(path string, parent *File) (*File, error) {
	return newFile(path, parent)
}

func newFile(path string, parent *File) (*File, error) {

	route, err := route.New(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a File for the path %q. Error: %s", path, err)
	}

	return &File{
		route:  route,
		parent: parent,
	}, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.route)
}

func (file *File) Route() *route.Route {
	return file.route
}

func (file *File) SetParent(parent *File) {
	file.parent = parent
}

func (file *File) Parent() *File {
	return file.parent
}

func (file *File) SetChilds(childs []*File) {

	// make the the current File the parent for all childs
	for _, child := range childs {
		child.SetParent(file)
	}

	file.childs = childs
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
