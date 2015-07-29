// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"strings"
	"text/template"

	"allmark.io/modules/model"
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

	AliasIndexTemplateName        = "aliasindex"
	AliasIndexContentTemplateName = "aliasindexcontent"

	SearchTemplateName        = "search"
	SearchContentTemplateName = "searchcontent"

	ConversionTemplateName = "converter"

	AliasesSnippetTemplateName   = "aliases-snippet"
	TagsSnippetTemplateName      = "tags-snippet"
	PublisherSnippetTemplateName = "publisher-snippet"

	// RobotsTxtTemplateName defines the name of the robots.txt template.
	RobotsTxtTemplateName = "robotstxt"
)

type Provider struct {
	Modified chan bool

	folder    string
	templates map[string]*Template
	snippets  map[string]*Template
	cache     map[string]*template.Template
}

func NewProvider(templateFolder string) Provider {

	// register all templates
	templates := make(map[string]*Template)
	templates[MasterTemplateName] = NewTemplate(templateFolder, MasterTemplateName, masterTemplate)
	templates[ErrorTemplateName] = NewTemplate(templateFolder, ErrorTemplateName, errorTemplate)
	templates[OpenSearchDescriptionTemplateName] = NewTemplate(templateFolder, OpenSearchDescriptionTemplateName, openSearchDescriptionTemplate)
	templates[TagmapTemplateName] = NewTemplate(templateFolder, TagmapTemplateName, tagmapTemplate)
	templates[TagmapContentTemplateName] = NewTemplate(templateFolder, TagmapContentTemplateName, tagmapContentTemplate)
	templates[AliasIndexTemplateName] = NewTemplate(templateFolder, AliasIndexTemplateName, aliasIndexTemplate)
	templates[AliasIndexContentTemplateName] = NewTemplate(templateFolder, AliasIndexContentTemplateName, aliasIndexContentTemplate)
	templates[SitemapTemplateName] = NewTemplate(templateFolder, SitemapTemplateName, sitemapTemplate)
	templates[SitemapContentTemplateName] = NewTemplate(templateFolder, SitemapContentTemplateName, sitemapContentTemplate)
	templates[XmlSitemapTemplateName] = NewTemplate(templateFolder, XmlSitemapTemplateName, xmlSitemapTemplate)
	templates[XmlSitemapContentTemplateName] = NewTemplate(templateFolder, XmlSitemapContentTemplateName, xmlSitemapContentTemplate)
	templates[RssFeedTemplateName] = NewTemplate(templateFolder, RssFeedTemplateName, rssFeedTemplate)
	templates[RssFeedContentTemplateName] = NewTemplate(templateFolder, RssFeedContentTemplateName, rssFeedContentTemplate)
	templates[SearchTemplateName] = NewTemplate(templateFolder, SearchTemplateName, searchTemplate)
	templates[SearchContentTemplateName] = NewTemplate(templateFolder, SearchContentTemplateName, searchContentTemplate)
	templates[model.TypeDocument.String()] = NewTemplate(templateFolder, model.TypeDocument.String(), documentTemplate)
	templates[model.TypeRepository.String()] = NewTemplate(templateFolder, model.TypeRepository.String(), repositoryTemplate)
	templates[model.TypePresentation.String()] = NewTemplate(templateFolder, model.TypePresentation.String(), presentationTemplate)
	templates[ConversionTemplateName] = NewTemplate(templateFolder, ConversionTemplateName, converterTemplate)
	templates[RobotsTxtTemplateName] = NewTemplate(templateFolder, RobotsTxtTemplateName, robotsTxtTemplate)

	// register snippets
	snippets := make(map[string]*Template)
	snippets[AliasesSnippetTemplateName] = NewTemplate(templateFolder, AliasesSnippetTemplateName, aliasesSnippet)
	snippets[TagsSnippetTemplateName] = NewTemplate(templateFolder, TagsSnippetTemplateName, tagsSnippet)
	snippets[PublisherSnippetTemplateName] = NewTemplate(templateFolder, PublisherSnippetTemplateName, publisherSnippet)

	// create the provider
	provider := Provider{
		folder:    templateFolder,
		templates: templates,
		snippets:  snippets,
		cache:     make(map[string]*template.Template),
	}

	return provider
}

func (provider *Provider) GetFullTemplate(hostname, templateName string) (*template.Template, error) {

	t, err := provider.getParsedTemplate(templateName, true)
	if err != nil {
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
	getAbsoluteURL := func(uri string) string {

		// sanatize
		uri = strings.TrimSpace(uri)

		// add prefix
		if !strings.HasPrefix(uri, "/") {
			uri = "/" + uri
		}

		return getHostname() + uri
	}

	// Replace all occurances of `textToReplace` in `text` with `replacement`.
	replace := func(text, textToReplace, replacement string) string {
		return strings.Replace(text, textToReplace, replacement, -1)
	}

	return map[string]interface{}{
		"hostname": getHostname,
		"absolute": getAbsoluteURL,
		"replace":  replace,
	}
}

func (provider *Provider) StoreTemplatesOnDisc() (success bool, err error) {

	// store snippets on disk
	for _, template := range provider.snippets {
		if savedToDisc, err := template.StoreOnDisc(); !savedToDisc {
			return false, err
		}
	}

	// store templates on disk
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

	// parse the snippets
	for _, snippet := range provider.snippets {

		template, err = template.Parse(snippet.Text())
		if err != nil {
			return nil, fmt.Errorf("Error while parsing snippet %q. Error: %s", snippet.Name(), err.Error())
		}

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
