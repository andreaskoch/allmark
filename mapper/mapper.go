package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
)

//type Mapper func(item *indexer.Item, childItemCallback func(item *indexer.Item)) interface{}

func GetMapper(item *indexer.Item) (func(item *indexer.Item, childItemCallback func(item *indexer.Item)) interface{}, error) {

	switch item.Type {
	case indexer.DocumentItemType:
		return documentMapperFunc, nil

	case indexer.MessageItemType:
		return messageMapperFunc, nil

	case indexer.CollectionItemType:
		return collectionMapperFunc, nil

	case indexer.RepositoryItemType:
		return repositoryMapperFunc, nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))

}
