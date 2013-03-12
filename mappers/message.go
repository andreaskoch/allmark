package mappers

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func GetMessage(item indexer.Item) viewmodel.Message {
	return viewmodel.Message{
		Title:       getTitle(item),
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getTitle(item indexer.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
