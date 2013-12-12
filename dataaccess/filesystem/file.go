// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"io/ioutil"
	"path/filepath"
)

func getFiles(repository *Repository, directory string) []*dataaccess.File {

	childs := make([]*dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		path := filepath.Join(directory, directoryEntry.Name())

		// recurse if the path is a directory
		if isDir, _ := fsutil.IsDirectory(path); isDir {
			childs = append(childs, getFiles(repository, path)...)
			continue
		}

		// append new file
		file, err := newFile(repository, path)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", path, err)
		}

		childs = append(childs, file)
	}

	return childs
}

func newFile(repository *Repository, path string) (*dataaccess.File, error) {

	// check if the path is a file
	if isFile, _ := fsutil.IsFile(path); !isFile {
		return nil, fmt.Errorf("%q is not a file.", path)
	}

	// route
	route, err := route.NewFromPath(repository.Path(), path)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a File for the path %q. Error: %s", path, err)
	}

	// content provider
	contentProvider := newContentProvider(path, route)

	// create the file
	file, err := dataaccess.NewFile(route, contentProvider)

	if err != nil {
		return nil, err
	}

	return file, nil
}
