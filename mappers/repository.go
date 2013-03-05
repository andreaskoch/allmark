package mappers

import (
	"fmt"
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
	"path/filepath"
)

func GetRepository(item indexer.Item) viewmodel.Repository {

	itemDirectory := filepath.Dir(item.Path)

	fmt.Println(itemDirectory)

	entries := make([]viewmodel.RepositoryEntry, 0, len(item.ChildItems))
	for _, child := range item.ChildItems {
		newEntry := viewmodel.RepositoryEntry{Path: child.GetRelativePath(itemDirectory)}
		entries = append(entries, newEntry)
	}

	return viewmodel.Repository{
		Title:       item.GetBlockValue("title"),
		Description: item.GetBlockValue("description"),
		Content:     item.GetBlockValue("content"),
		Entries:     entries,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
