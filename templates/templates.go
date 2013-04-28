// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/andreaskoch/allmark/parser"
	"strings"
)

const (
	ChildTemplatePlaceholder = "@childtemplate"
)

func GetTemplate(itemType string) string {
	masterTemplate := getMasterTemplate()
	childTempalte := getChildTemplate(itemType)

	return strings.Replace(masterTemplate, ChildTemplatePlaceholder, childTempalte, 1)
}

func getMasterTemplate() string {
	return masterTemplate
}

func getChildTemplate(itemType string) string {

	switch itemType {
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
