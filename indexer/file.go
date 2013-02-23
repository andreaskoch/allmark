package indexer

type File struct {
	Path string
}

func NewFile(path string) File {
	return File{
		Path: path,
	}
}
