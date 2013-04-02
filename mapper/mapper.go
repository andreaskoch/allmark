package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func GetMapper(item *repository.Item) (func(item *repository.Item) view.Model, error) {

	switch item.Type {
	case repository.DocumentItemType, repository.CollectionItemType:
		return documentMapperFunc, nil

	case repository.MessageItemType:
		return messageMapperFunc, nil

	case repository.RepositoryItemType:
		return repositoryMapperFunc, nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))

}
