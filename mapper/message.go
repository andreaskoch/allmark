package mapper

import (
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func createMessageMapperFunc(pathProvider *path.Provider) func(item *repository.Item) view.Model {
	return func(item *repository.Item) view.Model {
		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       getTitle(item),
			Content:     item.Content,
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}
}

func getTitle(item *repository.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
