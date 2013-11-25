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
	route *route.Route

	hashProvider    HashProviderFunc
	contentProvider ContentProviderFunc
}

func NewFile(path string, hashProvider HashProviderFunc, contentProvider ContentProviderFunc) (*File, error) {

	route, err := route.New(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a File for the path %q. Error: %s", path, err)
	}

	return &File{
		route: route,

		hashProvider:    hashProvider,
		contentProvider: contentProvider,
	}, nil
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.route)
}

func (file *File) Route() *route.Route {
	return file.route
}

func (file *File) HashProvider() HashProviderFunc {
	return file.hashProvider
}

func (file *File) ContentProvider() ContentProviderFunc {
	return file.contentProvider
}
