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
