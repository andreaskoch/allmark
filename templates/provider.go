// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark/types"
	"strings"
	"text/template"
)

const (
	// Template placholders
	ChildTemplatePlaceholder = "@childtemplate"

	// Template names
	MasterTemplateName = "master"
	ErrorTemplateName  = "error"
)

type Provider struct {
	Modified chan bool

	folder    string
	templates map[string]*Template
	cache     map[string]*template.Template
}

func NewProvider(templateFolder string) *Provider {

	// intialize the templates
	templateModified := make(chan bool)
	templates := make(map[string]*Template)

	templates[MasterTemplateName] = NewTemplate(templateFolder, MasterTemplateName, masterTemplate, templateModified)
	templates[ErrorTemplateName] = NewTemplate(templateFolder, ErrorTemplateName, errorTemplate, templateModified)

	templates[types.DocumentItemType] = NewTemplate(templateFolder, types.DocumentItemType, documentTemplate, templateModified)
	templates[types.MessageItemType] = NewTemplate(templateFolder, types.MessageItemType, messageTemplate, templateModified)
	templates[types.RepositoryItemType] = NewTemplate(templateFolder, types.RepositoryItemType, repositoryTemplate, templateModified)
	templates[types.PresentationItemType] = NewTemplate(templateFolder, types.PresentationItemType, presentationTemplate, templateModified)

	// create the provider
	provider := &Provider{
		Modified: make(chan bool),

		folder:    templateFolder,
		templates: templates,
		cache:     make(map[string]*template.Template),
	}

	// watch for changes
	go func() {
		for {
			select {
			case <-templateModified:
				provider.ClearCache()
				go func() {
					provider.Modified <- true
				}()
			}
		}
	}()

	return provider
}

func (provider *Provider) GetTemplate(itemType string) (*template.Template, error) {

	// get template from cache
	if template, ok := provider.cache[itemType]; ok {
		return template, nil
	}

	// assemble to the template
	templateText, err := provider.getTemplateText(itemType)
	if err != nil {
		return nil, err
	}

	// parse the template
	template, err := template.New(itemType).Parse(templateText)
	if err != nil {
		return nil, err
	}

	// add template to cache
	provider.cache[itemType] = template

	return template, nil
}

func (provider *Provider) StoreTemplatesOnDisc() (success bool, err error) {
	for _, template := range provider.templates {
		if savedToDisc, err := template.StoreOnDisc(); !savedToDisc {
			return false, err
		}
	}

	return true, nil
}

func (provider *Provider) ClearCache() {
	fmt.Println("Clearing the template cache.")

	provider.cache = make(map[string]*template.Template)
}

func (provider *Provider) getTemplateText(childTemplateName string) (string, error) {

	// get the master template
	masterTemplate := provider.getTemplate(MasterTemplateName)
	if masterTemplate == nil {
		return "", fmt.Errorf("Master template not found.")
	}

	// get the child template
	childTemplate := provider.getTemplate(childTemplateName)
	if childTemplate == nil {
		return "", fmt.Errorf("Child template %q not found.", childTemplateName)
	}

	// merge master and child template
	mergedTemplate := strings.Replace(masterTemplate.Text(), ChildTemplatePlaceholder, childTemplate.Text(), 1)

	return mergedTemplate, nil
}

func (provider *Provider) getTemplate(templateName string) *Template {

	if template, exists := provider.templates[templateName]; exists {
		return template
	}

	return nil
}
