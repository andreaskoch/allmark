package mapper

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func documentMapperFunc(item *indexer.Item, childItemCallback func(item *indexer.Item)) interface{} {
	return viewmodel.Document{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
