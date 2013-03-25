package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/viewmodel"
	"path/filepath"
)

func collectionMapperFunc(item *repository.Item, childItemCallback func(item *repository.Item)) interface{} {

	return viewmodel.Collection{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Entries:     getCollectionEntries(item, childItemCallback),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getCollectionEntries(item *repository.Item, childItemCallback func(item *repository.Item)) []viewmodel.CollectionEntry {
	parentDirectory := filepath.Dir(item.Path)

	getCollectionEntry := func(item repository.Item) viewmodel.CollectionEntry {

		return viewmodel.CollectionEntry{
			Title: item.Title,
			Path:  item.GetRelativePath(parentDirectory),
		}
	}

	entries := make([]viewmodel.CollectionEntry, 0, len(item.ChildItems))
	for _, child := range item.ChildItems {
		childItemCallback(child)
		entries = append(entries, getCollectionEntry(*child))
	}

	return entries
}
