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

	directoryName := strings.ToLower(filepath.Base(path))
	for _, reservedDirectoryName := range ReservedDirectoryNames {
		if directoryName == strings.ToLower(reservedDirectoryName) {
			return true
		}
	}

	return false
}
