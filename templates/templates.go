// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/andreaskoch/allmark/parser"
)

func GetTemplate(parserResult *parser.Result) string {

	switch itemType := parserResult.MetaData.ItemType; itemType {
	case parser.DocumentItemType:
		return documentTemplate

	case parser.MessageItemType:
		return messageTemplate

	case parser.CollectionItemType:
		return collectionTemplate

	case parser.RepositoryItemType:
		return repositoryTemplate

	default:
		return errorTemplate
	}

	panic("Unreachable")

}
