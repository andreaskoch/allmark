package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func GetDocument(item indexer.Item) viewmodel.Document {
	return viewmodel.Document{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.GetBlockValue("content"),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
