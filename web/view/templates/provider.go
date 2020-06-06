// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/elWyatt/allmark/web/view/templates/defaulttheme"
	"github.com/elWyatt/allmark/web/view/templates/templatenames"
	"github.com/elWyatt/allmark/web/webpaths"
)

// A Provider gives access to all required templates.
type Provider struct {
	Modified chan bool

	folder              string
	templatedefinitions map[string]*templateDefinition
}

// NewProvider creates a new template provider with the given folder as the base.
func NewProvider(templateFolder string) Provider {

	// register all templates
	templates := make(map[string]*templateDefinition)
	for templateName, rawTemplate := range defaulttheme.RawTemplates() {
		templates[templateName] = newTemplateDefinition(templateFolder, templateName, rawTemplate)
	}

	// create the provider
	provider := Provider{
		folder:              templateFolder,
		templatedefinitions: templates,
	}

	return provider
}

// GetErrorTemplate returns the template for error pages.
func (provider *Provider) GetErrorTemplate(hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(templatenames.Error, hostname)
}

// GetAliasIndexTemplate returns the alias-index template.
func (provider *Provider) GetAliasIndexTemplate(hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(templatenames.AliasIndex, hostname)
}

// GetSitemapTemplate returns the sitemap template.
func (provider *Provider) GetSitemapTemplate(hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(templatenames.Sitemap, hostname)
}

// GetSitemapEntryTemplate returns the sitemap-entry template.
func (provider *Provider) GetSitemapEntryTemplate(hostname string) (template *template.Template, childPlaceholder string, err error) {
	template, err = provider.GetSimpleTemplate(templatenames.SitemapEntry, hostname)
	return template, defaulttheme.SitemapChildPlaceholder, err
}

// GetSearchTemplate returns the search template.
func (provider *Provider) GetSearchTemplate(hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(templatenames.Search, hostname)
}

// GetItemTemplate returns the item template for the given item type (e.g. document, presentation).
func (provider *Provider) GetItemTemplate(itemType, hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(itemType, hostname)
}

// GetTagMapTemplate returns the template for tags.
func (provider *Provider) GetTagMapTemplate(hostname string) (*template.Template, error) {
	return provider.getWrappedTemplate(templatenames.TagMap, hostname)
}

// GetRSSTemplate returns the template for RSS feeds.
func (provider *Provider) GetRSSTemplate(hostname string) (*template.Template, error) {
	return provider.GetSimpleTemplate(templatenames.RSSFeed, hostname)
}

// GetXMLSitemapTemplate returns the template for XML sitemaps.
func (provider *Provider) GetXMLSitemapTemplate(hostname string) (*template.Template, error) {
	return provider.GetSimpleTemplate(templatenames.XMLSitemap, hostname)
}

// GetRobotsTxtTemplate returns the template for robots.txt.
func (provider *Provider) GetRobotsTxtTemplate(hostname string) (*template.Template, error) {
	return provider.GetSimpleTemplate(templatenames.RobotsTxt, hostname)
}

// GetConversionTemplate returns the template for conversion.
func (provider *Provider) GetConversionTemplate(hostname string) (*template.Template, error) {
	return provider.GetSimpleTemplate(templatenames.Conversion, hostname)
}

// GetOpenSearchDescriptionTemplate returns the template for conversion.
func (provider *Provider) GetOpenSearchDescriptionTemplate(hostname string) (*template.Template, error) {
	return provider.GetSimpleTemplate(templatenames.OpenSearchDescription, hostname)
}

// GetSimpleTemplate returns a simple template without wrapping or combination with other templates.
func (provider *Provider) GetSimpleTemplate(templateName, hostname string) (*template.Template, error) {

	// get the template code
	code, err := provider.getTemplateText(templateName)
	if err != nil {
		return nil, err
	}

	tmpl, err := provider.createTemplate(
		templateName,
		code,
		hostname)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// GetSnippetTemplate returns a snippet template without wrapping or combination with other templates.
func (provider *Provider) GetSnippetTemplate(snippetName, hostname string) (*template.Template, error) {

	// get the template code
	code, err := provider.getTemplateText(snippetName)
	if err != nil {
		return nil, err
	}

	code = fmt.Sprintf(`{{template "%s" .}}\n%s`, snippetName, code)

	tmpl, err := provider.createTemplate(
		snippetName,
		code,
		hostname)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// StoreTemplatesOnDisc saves all templates to disc.
func (provider *Provider) StoreTemplatesOnDisc() (success bool, err error) {

	// store templates definitions on disk
	for _, template := range provider.templatedefinitions {
		if savedToDisc, err := template.StoreOnDisc(); !savedToDisc {
			return false, err
		}
	}

	return true, nil
}

// getWrappedTemplate returns the supplied template wrapped by the master template
func (provider *Provider) getWrappedTemplate(subTemplate, hostname string) (*template.Template, error) {

	// get the master template code
	masterTemplateCode, err := provider.getTemplateText(templatenames.Master)
	if err != nil {
		return nil, err
	}

	// get the sub-template code
	subTemplateCode, err := provider.getTemplateText(subTemplate)
	if err != nil {
		return nil, err
	}

	// wrap the sub-template
	subTemplateCode = fmt.Sprintf(`{{define "content"}}%s{{end}}`, subTemplateCode)

	masterTemplate, err := provider.createTemplate(
		subTemplate,
		masterTemplateCode+subTemplateCode,
		hostname)

	if err != nil {
		return nil, err
	}

	return masterTemplate, nil
}

// createTemplate creates a template from the lateName, templateCode, hostname string) (*template.Template, error) {
func (provider *Provider) createTemplate(templateName, templateCode, hostname string) (*template.Template, error) {
	tmpl := template.Template{}
	tmpl.New(templateName).Funcs(getTemplateHelpers(hostname))

	// parse the template text
	_, err := tmpl.Parse(templateCode)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template %q. Error: %s", templateName, err.Error())
	}

	return &tmpl, nil
}

// getTemplateText returns the text of the template with the given name
func (provider *Provider) getTemplateText(templateName string) (string, error) {

	if template, exists := provider.templatedefinitions[templateName]; exists {
		return template.Text(), nil
	}

	return "", fmt.Errorf("The template with the name %q was not found.", templateName)
}

// getTemplateHelpers returns a map of utility functions that can be used in the templates.
func getTemplateHelpers(hostname string) map[string]interface{} {

	// Get the current hostname
	getHostname := func() string {
		return hostname
	}

	// get the absolute url for a given (relative) uri
	getAbsoluteURL := func(uri string) string {

		if webpaths.IsAbsoluteURI(uri) {
			return uri
		}

		// sanatize
		uri = strings.TrimSpace(uri)

		// add prefix
		if !strings.HasPrefix(uri, "/") {
			uri = "/" + uri
		}

		return getHostname() + uri
	}

	return map[string]interface{}{
		"hostname": getHostname,
		"absolute": getAbsoluteURL,
		"replace":  replace,
	}
}

// Replace all occurances of `textToReplace` in `text` with `replacement`.
func replace(text, textToReplace, replacement string) string {
	return strings.Replace(text, textToReplace, replacement, -1)
}
