// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
	"strings"
)

var (
	ReservedDirectoryNames = []string{FilesDirectoryName, config.MetaDataFolderName}
)

func isReservedDirectory(path string) bool {
	if isDirectory, _ := util.IsDirectory(path); !isDirectory {
		return false
	}

	// get the directory name
	directoryName := strings.ToLower(filepath.Base(path))

	// all dot-directories are ignored
	if strings.HasPrefix(directoryName, ".") {
		return true
	}

	// check the reserved directory names
	for _, reservedDirectoryName := range ReservedDirectoryNames {
		if directoryName == strings.ToLower(reservedDirectoryName) {
			return true
		}
	}

	return false
}
