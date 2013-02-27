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

func (index Index) GetRelativeItemPaths() []string {

	paths := make([]string, 0, 0)
	for _, item := range index.Items {
		paths = append(paths, item.GetRelativeItemPaths(index.Path)...)
	}

	return paths
}
