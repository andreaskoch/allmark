// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/andreaskoch/allmark/parser"
	"strings"
	"text/template"
)

const (
	ChildTemplatePlaceholder = "@childtemplate"
)

type TemplateProvider struct {
	folder string
	cache  map[string]*template.Template
}

func New(templateFolder string) *TemplateProvider {
	return &TemplateProvider{
		folder: templateFolder,
		cache:  make(map[string]*template.Template),
	}
}

func (templateProvider *TemplateProvider) GetTemplate(itemType string) (*template.Template, error) {

	// get template from cache
	if template, ok := templateProvider.cache[itemType]; ok {
		return template, nil
	}

	// assemble to the template
	masterTemplate := templateProvider.getMasterTemplate()
	childTempalte := templateProvider.getChildTemplate(itemType)
	mergedTemplateText := strings.Replace(masterTemplate, ChildTemplatePlaceholder, childTempalte, 1)

	// parse the template
	template, err := template.New(itemType).Parse(mergedTemplateText)
	if err != nil {
		return nil, err
	}

	// add template to cache
	templateProvider.cache[itemType] = template

	return template, nil
}

func (templateProvider *TemplateProvider) getMasterTemplate() string {
	return masterTemplate
}

func (templateProvider *TemplateProvider) getChildTemplate(itemType string) string {

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
