package model

import "fmt"

type RepositoryIndex struct {
	Path  string
	Items []RepositoryItem
}

func NewRepositoryIndex(path string, items []RepositoryItem) RepositoryIndex {
	return RepositoryIndex{
		Path:  path,
		Items: items,
	}
}

func (index *RepositoryIndex) ToString() string {
	s := ""

	for index, repositoryItem := range index.Items {
		s += fmt.Sprintf("%v)\n%v\n", index, repositoryItem.String())
	}

	return s
}
