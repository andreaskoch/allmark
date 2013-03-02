package indexer

type Index struct {
	Path  string
	items []Item
}

func NewIndex(path string, items []Item) Index {
	return Index{
		Path:  path,
		items: items,
	}
}

func (index Index) GetAllItems() []Item {

	items := make([]Item, 0, 0)

	index.Walk(func(item Item) {
		items = append(items, item)
	})

	return items
}

func (index Index) Walk(walkFunc func(item Item)) {
	for _, item := range index.items {
		item.Walk(walkFunc)
	}
}
