package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/viewmodel"
	"path/filepath"
)

func repositoryMapperFunc(item *repository.Item, childItemCallback func(item *repository.Item)) interface{} {

	return viewmodel.Repository{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Entries:     getRepositoryEntries(item, childItemCallback),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getRepositoryEntries(item *repository.Item, childItemCallback func(item *repository.Item)) []viewmodel.RepositoryEntry {
	parentDirectory := filepath.Dir(item.Path)

	getRepositoryEntry := func(item repository.Item) viewmodel.RepositoryEntry {

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
