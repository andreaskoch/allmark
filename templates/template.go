// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"errors"
	"fmt"
	"github.com/andreaskoch/allmark/repository"
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
