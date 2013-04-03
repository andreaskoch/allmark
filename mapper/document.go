package mapper

import (
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func createDocumentMapperFunc(pathProvider *path.Provider) func(item *repository.Item) view.Model {

	return func(item *repository.Item) view.Model {
		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       item.Title,
			Description: item.Description,
			Content:     item.Content,
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}

}
