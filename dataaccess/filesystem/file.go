// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"io/ioutil"
	"path/filepath"
)

func getFiles(repository *Repository, itemDirectory string) []*dataaccess.File {

	// get the "files"-directory
	filesDirectory := filepath.Join(itemDirectory, config.FilesDirectoryName)

	childs := make([]*dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(filesDirectory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		filePath := filepath.Join(filesDirectory, directoryEntry.Name())

		// recurse if the path is a directory
		if isDir, _ := fsutil.IsDirectory(filePath); isDir {
			childs = append(childs, getFiles(repository, filePath)...)
			continue
		}

		// append new file
		file, err := newFile(itemDirectory, filePath)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", filePath, err)
		}

		childs = append(childs, file)
	}

	return childs
}

func newFile(basePath, filePath string) (*dataaccess.File, error) {

	// check if the file path is a file
	if isFile, _ := fsutil.IsFile(filePath); !isFile {
		return nil, fmt.Errorf("%q is not a file.", filePath)
	}

	// route
	route, err := route.NewFromFilePath(basePath, filePath)
	if err != nil {
		return nil, fmt.Errorf("Cannot create a File for the file path %q. Error: %s", filePath, err)
	}

	// content provider
	contentProvider := newContentProvider(filePath, route)

	// create the file
	file, err := dataaccess.NewFile(route, contentProvider)

	if err != nil {
		return nil, err
	}

	return file, nil
}
