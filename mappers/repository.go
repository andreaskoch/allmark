package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
	"path/filepath"
)

func GetRepository(item indexer.Item, childItemCallback func(item *indexer.Item)) viewmodel.Repository {

	return viewmodel.Repository{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Entries:     getRepositoryEntries(item, childItemCallback),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getRepositoryEntries(item indexer.Item, childItemCallback func(item *indexer.Item)) []viewmodel.RepositoryEntry {
	parentDirectory := filepath.Dir(item.Path)

	getRepositoryEntry := func(item indexer.Item) viewmodel.RepositoryEntry {

		return viewmodel.RepositoryEntry{
			Title: item.Title,
			Path:  item.GetRelativePath(parentDirectory),
		}
	}

	entries := make([]viewmodel.RepositoryEntry, 0, len(item.ChildItems))
	for _, child := range item.ChildItems {
		childItemCallback(child)
		entries = append(entries, getRepositoryEntry(*child))
	}

	return entries
}
