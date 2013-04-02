package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func documentMapperFunc(item *repository.Item, pathProviderFunc func(item *repository.Item) string) view.Model {
	return view.Model{
		Path:        pathProviderFunc(item),
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
