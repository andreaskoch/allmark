package mappers

import (
	"andyk/docs/indexer"
	"andyk/docs/viewmodel"
)

func GetDocument(item indexer.Item) viewmodel.Document {
	return viewmodel.Document{
		Title:       item.GetBlockValue("title"),
		Description: item.GetBlockValue("description"),
		Content:     item.GetBlockValue("content"),
		//LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
	}
}
