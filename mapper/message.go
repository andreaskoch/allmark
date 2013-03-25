package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/viewmodel"
)

func messageMapperFunc(item *repository.Item, childItemCallback func(item *repository.Item)) interface{} {
	return viewmodel.Message{
		Title:       getTitle(item),
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getTitle(item *repository.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
