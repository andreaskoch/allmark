package repository

import "fmt"

type Index struct {
	Path  string
	Items []Item
}

func NewIndex(path string, items []Item) Index {
	return Index{
		Path:  path,
		Items: items,
	}
}

func (index *Index) ToString() string {
	s := ""

	for index, item := range index.Items {
		s += fmt.Sprintf("%v)\n%v\n", index, item.String())
	}

	return s
}
