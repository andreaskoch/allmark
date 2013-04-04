package repository

type File struct {
	path string
}

func NewFile(filePath string) *File {
	return &File{
		path: filePath,
	}
}

func (file *File) Path() string {
	return file.path
}
func (file *File) PathType() string {
	return "file"
}
