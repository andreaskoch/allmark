package mapper

import (
	"fmt"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func repositoryMapperFunc(item *repository.Item, pathProviderFunc func(item *repository.Item) string) view.Model {

	return view.Model{
		Path:        pathProviderFunc(item),
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		Entries:     getEntries(item, pathProviderFunc),
		LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
	}
}

func getEntries(item *repository.Item, pathProviderFunc func(item *repository.Item) string) []view.Model {

	entries := make([]view.Model, 0)

	for _, child := range item.ChildItems {
		if mapperFunc, err := GetMapper(child, pathProviderFunc); err == nil {
			viewModel := mapperFunc(child)
			entries = append(entries, viewModel)
		} else {
			fmt.Println(err)
		}

	}

	return entries
}
