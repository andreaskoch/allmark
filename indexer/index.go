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
	minSize := len(index.Items) + 1
	items := make([]Item, minSize, minSize)

	// add all index items
	for _, item := range index.Items {
		items = append(items, item.GetAllItems()...)
	}

	return items
}
