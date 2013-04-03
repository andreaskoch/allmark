package repository

type Index struct {
	path  string
	items []*Item
}

func NewIndex(indexPath string, items []*Item) *Index {
	return &Index{
		path:  indexPath,
		items: items,
	}
}

func (index *Index) Walk(walkFunc func(item *Item)) {
	for _, item := range index.items {
		item.Walk(walkFunc)
	}
}

func (index *Index) DirectoryAbsolute() string {
	return index.path
}
