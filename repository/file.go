package repository

import (
	"os"
	"strings"
)

type File struct {
	Path string
}

func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

func (file File) AbsolutePath() string {
	return file.Path
}

func (file File) RelativePath(basePath string) string {

	pathSeperator := string(os.PathSeparator)
	fullPath := file.Path
	relativePath := strings.Replace(fullPath, basePath, "", 1)
	relativePath = pathSeperator + strings.TrimLeft(relativePath, pathSeperator)
	return relativePath
}
