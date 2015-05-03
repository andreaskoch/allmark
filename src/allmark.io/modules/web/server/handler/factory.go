// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/templates"
)

func NewFactory(logger logger.Logger, config config.Config, templateProvider templates.Provider, orchestratorFactory orchestrator.Factory, headerWriterFactory header.WriterFactory) *Factory {
	return &Factory{
		logger: logger,
		config: config,

		templateProvider:    templateProvider,
		orchestratorFactory: orchestratorFactory,
		headerWriterFactory: headerWriterFactory,
	}
}

type Factory struct {
	logger logger.Logger
	config config.Config

	templateProvider    templates.Provider
	orchestratorFactory orchestrator.Factory
	headerWriterFactory header.WriterFactory
}

func (factory *Factory) NewErrorHandler() Handler {

	return &Error{
		logger:                 factory.logger,
		headerWriter:           factory.headerWriterFactory.Static(),
		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
	}

}

func (factory *Factory) NewRobotsTxtHandler() Handler {

	return &RobotsTxt{
		logger:       factory.logger,
		headerWriter: factory.headerWriterFactory.Static(),
	}

}

func (factory *Factory) NewXmlSitemapHandler() Handler {

	return &XmlSitemap{
		logger: factory.logger,

		headerWriter:           factory.headerWriterFactory.Static(),
		templateProvider:       factory.templateProvider,
		xmlSitemapOrchestrator: factory.orchestratorFactory.NewXmlSitemapOrchestrator(),
	}

}

func (factory *Factory) NewTagsHandler() Handler {

	return &Tags{
		logger: factory.logger,

		headerWriter:           factory.headerWriterFactory.Static(),
		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		tagsOrchestrator:       factory.orchestratorFactory.NewTagsOrchestrator(),
	}

}

func (factory *Factory) NewSitemapHandler() Handler {

	return &Sitemap{
		logger: factory.logger,

		headerWriter:           factory.headerWriterFactory.Static(),
		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		sitemapOrchestrator:    factory.orchestratorFactory.NewSitemapOrchestrator(),
	}

}

func (factory *Factory) NewSearchHandler() Handler {

	return &Search{
		logger: factory.logger,

		headerWriter:     factory.headerWriterFactory.Static(),
		templateProvider: factory.templateProvider,
		error404Handler:  factory.NewErrorHandler(),

		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		searchOrchestrator:     factory.orchestratorFactory.NewSearchOrchestrator(),
	}

}

func (factory *Factory) NewOpenSearchDescriptionHandler() Handler {

	return &OpenSearchDescription{
		logger: factory.logger,

		headerWriter:                      factory.headerWriterFactory.Static(),
		openSearchDescriptionOrchestrator: factory.orchestratorFactory.NewOpenSearchDescriptionOrchestrator(),
		templateProvider:                  factory.templateProvider,
	}

}

func (factory *Factory) NewTypeAheadSearchHandler() Handler {
	return &TypeAhead{
		logger: factory.logger,

		headerWriter:          factory.headerWriterFactory.Static(),
		typeAheadOrchestrator: factory.orchestratorFactory.NewTypeAheadOrchestrator(),
	}

}

func (factory *Factory) NewTypeAheadTitlesHandler() Handler {

	return &Titles{
		logger: factory.logger,

		headerWriter:       factory.headerWriterFactory.Static(),
		titlesOrchestrator: factory.orchestratorFactory.NewTitlesOrchestrator(),
	}

}

func (factory *Factory) NewRssHandler() Handler {

	return &Rss{
		logger: factory.logger,

		headerWriter:     factory.headerWriterFactory.Static(),
		templateProvider: factory.templateProvider,
		error404Handler:  factory.NewErrorHandler(),
		feedOrchestrator: factory.orchestratorFactory.NewFeedOrchestrator(),
	}

}

func (factory *Factory) NewItemHandler() Handler {

	return &Item{
		logger: factory.logger,

		headerWriter:          factory.headerWriterFactory.Static(),
		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fileOrchestrator:      factory.orchestratorFactory.NewFileOrchestrator(),
		templateProvider:      factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewPrintHandler() Handler {

	return &Print{

		logger: factory.logger,

		headerWriter:               factory.headerWriterFactory.Static(),
		converterModelOrchestrator: factory.orchestratorFactory.NewConversionModelOrchestrator(),
		templateProvider:           factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewRtfHandler() Handler {

	// check if rtf conversion is enabled
	if !factory.config.Conversion.Rtf.Enabled {
		return factory.NewErrorHandler()
	}

	// check if the a rtf conversion tool has been supplied
	conversionToolPath := factory.config.Conversion.Rtf.Tool
	if conversionToolPath == "" {
		return factory.NewErrorHandler()
	}

	return &Rtf{
		logger: factory.logger,

		conversionToolPath:         conversionToolPath,
		headerWriter:               factory.headerWriterFactory.Static(),
		converterModelOrchestrator: factory.orchestratorFactory.NewConversionModelOrchestrator(),
		templateProvider:           factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewJsonHandler() Handler {

	return &Json{
		logger: factory.logger,

		headerWriter:          factory.headerWriterFactory.Static(),
		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fallbackHandler:       factory.NewItemHandler(),
	}

}

func (factory *Factory) NewLatestHandler() Handler {

	return &Latest{
		logger: factory.logger,

		headerWriter:          factory.headerWriterFactory.Static(),
		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fallbackHandler:       factory.NewItemHandler(),
	}

}

func (factory *Factory) NewUpdateHandler() *Update {

	return &Update{
		logger: factory.logger,

		headerWriter:       factory.headerWriterFactory.Static(),
		updateOrchestrator: factory.orchestratorFactory.NewUpdateOrchestrator(),
	}

}

func (factory *Factory) NewThemeHandler() Handler {

	return &Theme{
		logger: factory.logger,

		headerWriter:    factory.headerWriterFactory.Static(),
		error404Handler: factory.NewErrorHandler(),
	}

}
