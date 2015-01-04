// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"strings"
	"text/template"
)

const (
	// Template placholders
	ChildTemplatePlaceholder = "@childtemplate"

	// Template names
	MasterTemplateName = "master"
	ErrorTemplateName  = "error"

	OpenSearchDescriptionTemplateName = "opensearchdescription"

	SitemapTemplateName        = "sitemap"
	SitemapContentTemplateName = "sitemapcontent"

	XmlSitemapTemplateName        = "xmlsitemap"
	XmlSitemapContentTemplateName = "xmlsitemapcontent"

	RssFeedTemplateName        = "rssfeed"
	RssFeedContentTemplateName = "rssfeedcontent"

	TagmapTemplateName        = "tagmap"
	TagmapContentTemplateName = "tagmapcontent"

	SearchTemplateName        = "search"
	SearchContentTemplateName = "searchcontent"

	ConversionTemplateName = "converter"
)

type Provider struct {
	Modified chan bool

	folder    string
	templates map[string]*Template
	cache     map[string]*template.Template
}

func NewProvider(templateFolder string) Provider {

	// intialize the templates
	templates := make(map[string]*Template)

	templates[MasterTemplateName] = NewTemplate(templateFolder, MasterTemplateName, masterTemplate)
	templates[ErrorTemplateName] = NewTemplate(templateFolder, ErrorTemplateName, errorTemplate)

	templates[OpenSearchDescriptionTemplateName] = NewTemplate(templateFolder, OpenSearchDescriptionTemplateName, openSearchDescriptionTemplate)

	templates[TagmapTemplateName] = NewTemplate(templateFolder, TagmapTemplateName, tagmapTemplate)
	templates[TagmapContentTemplateName] = NewTemplate(templateFolder, TagmapContentTemplateName, tagmapContentTemplate)

	templates[SitemapTemplateName] = NewTemplate(templateFolder, SitemapTemplateName, sitemapTemplate)
	templates[SitemapContentTemplateName] = NewTemplate(templateFolder, SitemapContentTemplateName, sitemapContentTemplate)

	templates[XmlSitemapTemplateName] = NewTemplate(templateFolder, XmlSitemapTemplateName, xmlSitemapTemplate)
	templates[XmlSitemapContentTemplateName] = NewTemplate(templateFolder, XmlSitemapContentTemplateName, xmlSitemapContentTemplate)

	templates[RssFeedTemplateName] = NewTemplate(templateFolder, RssFeedTemplateName, rssFeedTemplate)
	templates[RssFeedContentTemplateName] = NewTemplate(templateFolder, RssFeedContentTemplateName, rssFeedContentTemplate)

	templates[SearchTemplateName] = NewTemplate(templateFolder, SearchTemplateName, searchTemplate)
	templates[SearchContentTemplateName] = NewTemplate(templateFolder, SearchContentTemplateName, searchContentTemplate)

	templates[model.TypeDocument.String()] = NewTemplate(templateFolder, model.TypeDocument.String(), documentTemplate)
	templates[model.TypeLocation.String()] = NewTemplate(templateFolder, model.TypeLocation.String(), locationTemplate)
	templates[model.TypeMessage.String()] = NewTemplate(templateFolder, model.TypeMessage.String(), messageTemplate)
	templates[model.TypeRepository.String()] = NewTemplate(templateFolder, model.TypeRepository.String(), repositoryTemplate)
	templates[model.TypePresentation.String()] = NewTemplate(templateFolder, model.TypePresentation.String(), presentationTemplate)

	templates[ConversionTemplateName] = NewTemplate(templateFolder, ConversionTemplateName, converterTemplate)

	// create the provider
	provider := Provider{
		folder:    templateFolder,
		templates: templates,
		cache:     make(map[string]*template.Template),
	}

	return provider
}

func (provider *Provider) GetFullTemplate(hostname, templateName string) (*template.Template, error) {

	t, err := provider.getParsedTemplate(templateName, true)
	if err != nil {
		panic(err)
		return nil, err
	}

	// override template functions
	t.Funcs(provider.getTemplateFunctions(hostname))

	return t, nil
}

func (provider *Provider) GetSubTemplate(hostname, templateName string) (*template.Template, error) {
	t, err := provider.getParsedTemplate(templateName, false)
	if err != nil {
		return nil, err
	}

	// override template functions
	t.Funcs(provider.getTemplateFunctions(hostname))

	return t, nil
}

func (provider *Provider) getTemplateFunctions(hostname string) map[string]interface{} {

	// Get the current hostname
	getHostname := func() string {
		return hostname
	}

	// get the absolute url for a given (relative) uri
	getAbsoluteUrl := func(uri string) string {

		// sanatize
		uri = strings.TrimSpace(uri)

		// add prefix
		if !strings.HasPrefix(uri, "/") {
			uri = "/" + uri
		}

		return "http://" + getHostname() + uri
	}

	// Replace all occurances of `textToReplace` in `text` with `replacement`.
	replace := func(text, textToReplace, replacement string) string {
		return strings.Replace(text, textToReplace, replacement, -1)
	}

	return map[string]interface{}{
		"hostname": getHostname,
		"absolute": getAbsoluteUrl,
		"replace":  replace,
	}
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
	// todo: use logger
	// fmt.Println("Clearing the template cache.")

	provider.cache = make(map[string]*template.Template)
}

func (provider *Provider) getParsedTemplate(templateName string, includeMaster bool) (*template.Template, error) {

	// get template from cache
	if template, ok := provider.cache[templateName]; ok {
		return template, nil
	}

	// assemble to the template
	childTemplate := provider.getTemplate(templateName)
	if childTemplate == nil {
		return nil, fmt.Errorf("Child template %q not found.", templateName)
	}

	templateText := childTemplate.Text()
	if includeMaster {

		masterTemplate := provider.getTemplate(MasterTemplateName)
		if masterTemplate == nil {
			return nil, fmt.Errorf("Master template not found.")
		}

		// merge master and child template
		templateText = strings.Replace(masterTemplate.Text(), ChildTemplatePlaceholder, templateText, 1)

	}

	// parse the template
	template, err := template.New(templateName).Funcs(provider.getTemplateFunctions("")).Parse(templateText)
	if err != nil {
		return nil, err
	}

	// add template to cache
	provider.cache[templateName] = template

	return template, nil
}

func (provider *Provider) getTemplate(templateName string) *Template {

	if template, exists := provider.templates[templateName]; exists {
		return template
	}

	return nil
}
