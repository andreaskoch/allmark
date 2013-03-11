package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
	"path/filepath"
)

func GetCollection(item indexer.Item, childItemCallback func(item *indexer.Item)) viewmodel.Collection {

	return viewmodel.Collection{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.GetBlockValue("content"),
		Entries:     getCollectionEntries(item, childItemCallback),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getCollectionEntries(item indexer.Item, childItemCallback func(item *indexer.Item)) []viewmodel.CollectionEntry {
	parentDirectory := filepath.Dir(item.Path)

	getCollectionEntry := func(item indexer.Item) viewmodel.CollectionEntry {

		return viewmodel.CollectionEntry{
			Title: item.Title,
			Path:  item.GetRelativePath(parentDirectory),
		}
	}

	entries := make([]viewmodel.CollectionEntry, 0, len(item.ChildItems))
	for _, child := range item.ChildItems {
		childItemCallback(&child)
		entries = append(entries, getCollectionEntry(child))
	}

	return entries
}
