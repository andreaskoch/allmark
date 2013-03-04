package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func GetDocument(item indexer.Item) viewmodel.Document {
	return viewmodel.Document{
		Title:       item.GetBlockValue("title"),
		Description: item.GetBlockValue("description"),
		Content:     item.GetBlockValue("content"),
		//LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
	}
}
