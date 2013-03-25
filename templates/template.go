package templates

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/repository"
)

func GetTemplate(item *repository.Item) (string, error) {

	switch itemType := item.Type; itemType {
	case repository.DocumentItemType:
		return documentTemplate, nil

	case repository.MessageItemType:
		return messageTemplate, nil

	case repository.CollectionItemType:
		return collectionTemplate, nil

	case repository.RepositoryItemType:
		return repositoryTemplate, nil

	default:
		return "", errors.New(fmt.Sprintf("No template available for items of type %q.", itemType))
	}

	panic("Unreachable")

}
