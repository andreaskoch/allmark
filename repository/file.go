package repository

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	path           string
	indexDirectory string
}

func NewFile(indexDirectory string, filePath string) *File {
	return &File{
		path:           filePath,
		indexDirectory: indexDirectory,
	}
}

func (file *File) IndexDirectoryAbsolute() string {
	return file.indexDirectory
}

func (file *File) DirectoryAbsolute() string {
	return filepath.Dir(file.path)
}

func (file *File) PathAbsolute() string {
	return file.path
}

func (file *File) Route() string {

	pathSeperator := string(os.PathSeparator)

	relativePath := strings.Replace(file.PathAbsolute(), file.IndexDirectoryAbsolute(), "", 1)
	relativePath = pathSeperator + strings.TrimLeft(relativePath, pathSeperator)
	relativePath = strings.Replace(relativePath, string(pathSeperator), "/", -1)

	return relativePath
}
