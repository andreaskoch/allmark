// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/andreaskoch/allmark/repository"
)

func GetTemplate(item *repository.Item) string {

	switch itemType := item.Type; itemType {
	case repository.DocumentItemType:
		return documentTemplate

	case repository.MessageItemType:
		return messageTemplate

	case repository.CollectionItemType:
		return collectionTemplate

	case repository.RepositoryItemType:
		return repositoryTemplate

	default:
		return errorTemplate
	}

	panic("Unreachable")

}
