package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/repository"
)

//type Mapper func(item *repository.Item, childItemCallback func(item *repository.Item)) interface{}

func GetMapper(item *repository.Item) (func(item *repository.Item, childItemCallback func(item *repository.Item)) interface{}, error) {

	switch item.Type {
	case repository.DocumentItemType:
		return documentMapperFunc, nil

	case repository.MessageItemType:
		return messageMapperFunc, nil

	case repository.CollectionItemType:
		return collectionMapperFunc, nil

	case repository.RepositoryItemType:
		return repositoryMapperFunc, nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))

}
