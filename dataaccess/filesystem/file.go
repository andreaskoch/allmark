// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"io/ioutil"
	"path/filepath"
)

func newRootFolder(path string) (*dataaccess.File, error) {
	return newFile(nil, path)
}

func newFile(parent *dataaccess.File, path string) (*dataaccess.File, error) {

	// check if the path exists
	if exists := fsutil.PathExists(path); !exists {
		return nil, fmt.Errorf("The path %q does not exists.", path)
	}

	// check if the path is a directory or a file
	isDir, err := fsutil.IsDirectory(path)
	if err != nil {
		return nil, err
	}

	// create the file
	file, err := dataaccess.NewFile(path, parent)
	if err != nil {
		return nil, err
	}

	// append childs
	if isDir {
		file.SetChilds(getChildFiles(path))
	}

	return file, nil
}

func getChildFiles(directory string) []*dataaccess.File {

	childs := make([]*dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		// append new file
		path := filepath.Join(directory, directoryEntry.Name())
		file, err := newRootFolder(path)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", path, err)
		}

		childs = append(childs, file)
	}

	return childs
}
