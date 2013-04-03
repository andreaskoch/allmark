package mapper

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/path"
	"github.com/andreaskoch/docs/repository"
	"github.com/andreaskoch/docs/view"
)

func GetMapper(pathProvider *path.Provider, item *repository.Item) (func(item *repository.Item) view.Model, error) {

	switch item.Type {
	case repository.DocumentItemType:
		return createDocumentMapperFunc(pathProvider), nil

	case repository.MessageItemType:
		return createMessageMapperFunc(pathProvider), nil

	case repository.RepositoryItemType, repository.CollectionItemType:
		return createCollectionMapperFunc(pathProvider), nil
	}

	return nil, errors.New(fmt.Sprintf("There is no mapper available for items of type %q", item.Type))
}
