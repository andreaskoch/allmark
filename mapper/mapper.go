package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func GetMapper(item *repository.Item, pathProviderFunc func(item *repository.Item) string) (func(item *repository.Item) view.Model, error) {

	switch item.Type {
	case repository.DocumentItemType, repository.CollectionItemType:

		return func(i *repository.Item) view.Model {
			return documentMapperFunc(i, pathProviderFunc)
		}, nil

	case repository.MessageItemType:

		return func(i *repository.Item) view.Model {
			return messageMapperFunc(i, pathProviderFunc)
		}, nil

	case repository.RepositoryItemType:

		return func(i *repository.Item) view.Model {
			return repositoryMapperFunc(i, pathProviderFunc)
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))

}
