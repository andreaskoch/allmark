package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func messageMapperFunc(item *repository.Item, pathProviderFunc func(item *repository.Item) string) view.Model {
	return view.Model{
		Path:        pathProviderFunc(item),
		Title:       getTitle(item),
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getTitle(item *repository.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
