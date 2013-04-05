// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"strings"
)

func AddLeadingFilesystemDirectorySeperator(filepath string) string {
	newPath := filepath

	for strings.Index(newPath, FilesystemDirectorySeperator) == 0 {
		newPath = strings.TrimLeft(newPath, FilesystemDirectorySeperator)
	}

	return FilesystemDirectorySeperator + newPath
}
