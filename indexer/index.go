package indexer

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

func (index Index) GetAllItems() []Item {

	// number of direct descendants plus the current item
	items := make([]Item, 0, 0)

	var walkFunc = func(item Item) {
		items = append(items, item)
	}

	// add all index items
	for _, item := range index.Items {
		item.Walk(walkFunc)
	}

	return items
}
