// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/view/templates"
)

func NewFactory(logger logger.Logger, config config.Config, templateProvider templates.Provider, orchestratorFactory orchestrator.Factory) *Factory {

	return &Factory{
		logger: logger,
		config: config,

		templateProvider:    templateProvider,
		orchestratorFactory: orchestratorFactory,
	}
}

type Factory struct {
	logger logger.Logger
	config config.Config

	templateProvider    templates.Provider
	orchestratorFactory orchestrator.Factory
}

func (factory *Factory) NewErrorHandler() Handler {

	return &Error{
		logger:                 factory.logger,
		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
	}

}

func (factory *Factory) NewRobotsTxtHandler() Handler {

	return &RobotsTxt{
		logger: factory.logger,
	}

}

func (factory *Factory) NewXmlSitemapHandler() Handler {

	return &XmlSitemap{
		logger: factory.logger,

		templateProvider:       factory.templateProvider,
		xmlSitemapOrchestrator: factory.orchestratorFactory.NewXmlSitemapOrchestrator(),
	}

}

func (factory *Factory) NewTagsHandler() Handler {

	return &Tags{
		logger: factory.logger,

		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		tagsOrchestrator:       factory.orchestratorFactory.NewTagsOrchestrator(),
	}

}

func (factory *Factory) NewSitemapHandler() Handler {

	return &Sitemap{
		logger: factory.logger,

		templateProvider:       factory.templateProvider,
		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		sitemapOrchestrator:    factory.orchestratorFactory.NewSitemapOrchestrator(),
	}

}

func (factory *Factory) NewSearchHandler() Handler {

	return &Search{
		logger: factory.logger,

		templateProvider: factory.templateProvider,
		error404Handler:  factory.NewErrorHandler(),

		navigationOrchestrator: factory.orchestratorFactory.NewNavigationOrchestrator(),
		searchOrchestrator:     factory.orchestratorFactory.NewSearchOrchestrator(),
	}

}

func (factory *Factory) NewOpenSearchDescriptionHandler() Handler {

	return &OpenSearchDescription{
		logger: factory.logger,

		openSearchDescriptionOrchestrator: factory.orchestratorFactory.NewOpenSearchDescriptionOrchestrator(),
		templateProvider:                  factory.templateProvider,
	}

}

func (factory *Factory) NewTypeAheadSearchHandler() Handler {
	return &TypeAhead{
		logger: factory.logger,

		typeAheadOrchestrator: factory.orchestratorFactory.NewTypeAheadOrchestrator(),
	}

}

func (factory *Factory) NewTypeAheadTitlesHandler() Handler {

	return &Titles{
		logger: factory.logger,

		titlesOrchestrator: factory.orchestratorFactory.NewTitlesOrchestrator(),
	}

}

func (factory *Factory) NewRssHandler() Handler {

	return &Rss{
		logger: factory.logger,

		templateProvider: factory.templateProvider,
		error404Handler:  factory.NewErrorHandler(),
		feedOrchestrator: factory.orchestratorFactory.NewFeedOrchestrator(),
	}

}

func (factory *Factory) NewItemHandler() Handler {

	return &Item{
		logger: factory.logger,

		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fileOrchestrator:      factory.orchestratorFactory.NewFileOrchestrator(),
		templateProvider:      factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewPrintHandler() Handler {

	return &Print{

		logger: factory.logger,

		converterModelOrchestrator: factory.orchestratorFactory.NewConversionModelOrchestrator(),
		templateProvider:           factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewRtfHandler() Handler {

	return &Rtf{
		logger: factory.logger,
		config: factory.config,

		converterModelOrchestrator: factory.orchestratorFactory.NewConversionModelOrchestrator(),
		templateProvider:           factory.templateProvider,

		error404Handler: factory.NewErrorHandler(),
	}

}

func (factory *Factory) NewJsonHandler() Handler {

	return &Json{
		logger: factory.logger,

		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fallbackHandler:       factory.NewItemHandler(),
	}

}

func (factory *Factory) NewLatestHandler() Handler {

	return &Latest{
		logger: factory.logger,

		viewModelOrchestrator: factory.orchestratorFactory.NewViewModelOrchestrator(),
		fallbackHandler:       factory.NewItemHandler(),
	}

}

func (factory *Factory) NewUpdateHandler() *Update {

	return &Update{
		logger: factory.logger,

		updateOrchestrator: factory.orchestratorFactory.NewUpdateOrchestrator(),
	}

}

func (factory *Factory) NewThemeHandler() Handler {

	return &Theme{
		logger: factory.logger,

		error404Handler: factory.NewErrorHandler(),
	}

}
