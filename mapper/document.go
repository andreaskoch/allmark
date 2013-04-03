package mapper

import (
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func documentMapperFunc(item *repository.Item) view.Model {
	return view.Model{
		Path:        item.Route(),
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}
