package mapper

import (
	"github.com/andreaskoch/docs/indexer"
	"github.com/andreaskoch/docs/viewmodel"
)

func messageMapperFunc(item *indexer.Item, childItemCallback func(item *indexer.Item)) interface{} {
	return viewmodel.Message{
		Title:       getTitle(item),
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getTitle(item *indexer.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
