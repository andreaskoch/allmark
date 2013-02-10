package model

type RepositoryItemFile struct {
	Path string
}

func NewRepositoryItemFile(path string) RepositoryItemFile {
	return RepositoryItemFile{
		Path: path,
	}
}
