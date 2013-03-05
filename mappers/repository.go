package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func GetRepository(item indexer.Item) viewmodel.Repository {
	return viewmodel.Repository{
		Title:       item.GetBlockValue("title"),
		Description: item.GetBlockValue("description"),
		Content:     item.GetBlockValue("content"),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
