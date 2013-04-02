package repository

import (
	"os"
	"strings"
)

type File struct {
	path string
}

func NewFile(filePath string) *File {
	return &File{
		path: filePath,
	}
}

func (file *File) PathAbsolute() string {
	return file.path
}

func (file *File) PathRelative() string {

	pathSeperator := string(os.PathSeparator)

	basePath := ""
	fullPath := file.path
	relativePath := strings.Replace(fullPath, basePath, "", 1)
	relativePath = pathSeperator + strings.TrimLeft(relativePath, pathSeperator)

	return relativePath
}
