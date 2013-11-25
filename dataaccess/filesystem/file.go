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

func getFiles(directory string) []*dataaccess.File {

	childs := make([]*dataaccess.File, 0)

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return childs
	}

	for _, directoryEntry := range filesDirectoryEntries {

		// append new file
		path := filepath.Join(directory, directoryEntry.Name())
		file, err := newFile(path)
		if err != nil {
			fmt.Printf("Unable to add file %q to index.\nError: %s\n", path, err)
		}

		childs = append(childs, file)
	}

	return childs
}

func newFile(path string) (*dataaccess.File, error) {

	// check if the path is a file
	if isFile, _ := fsutil.IsFile(path); !isFile {
		return nil, fmt.Errorf("%q is not a file.", path)
	}

	// hash provider
	hashProvider := func() (string, error) {
		return getHash(path)
	}

	// content provider
	contentProvider := func() ([]byte, error) {
		return getContent(path)
	}

	// create the file
	file, err := dataaccess.NewFile(path, hashProvider, contentProvider)

	if err != nil {
		return nil, err
	}

	return file, nil
}
