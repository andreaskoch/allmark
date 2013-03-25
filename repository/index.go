package repository

type Index struct {
	Path  string
	items []*Item
}

func NewIndex(path string, items []*Item) *Index {
	return &Index{
		Path:  path,
		items: items,
	}
}

func (index *Index) Walk(walkFunc func(item *Item)) {
	for _, item := range index.items {
		item.Walk(walkFunc)
	}
}
