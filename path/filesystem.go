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
