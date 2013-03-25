package templates

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/docs/indexer"
)

func GetTemplate(item *indexer.Item) (string, error) {

	switch itemType := item.Type; itemType {
	case indexer.DocumentItemType:
		return documentTemplate, nil

	case indexer.MessageItemType:
		return messageTemplate, nil

	case indexer.CollectionItemType:
		return collectionTemplate, nil

	case indexer.RepositoryItemType:
		return repositoryTemplate, nil

	default:
		return "", errors.New(fmt.Sprintf("No template available for items of type %q.", itemType))
	}

	panic("Unreachable")

}
