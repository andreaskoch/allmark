package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/viewmodel"
)

func documentMapperFunc(item *repository.Item, childItemCallback func(item *repository.Item)) interface{} {
	return viewmodel.Document{
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
