// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/webpaths"
)

func NewFactory(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Factory {

	baseOrchestrator := newBaseOrchestrator(logger, config, repository, parser, converter, webPathProvider)

	return &Factory{
		logger: logger,

		baseOrchestrator: baseOrchestrator,
	}
}

type Factory struct {
	logger logger.Logger

	baseOrchestrator *Orchestrator
}

func (factory *Factory) NewConversionModelOrchestrator() ConversionModelOrchestrator {
	return ConversionModelOrchestrator{
		Orchestrator:     factory.baseOrchestrator,
		fileOrchestrator: factory.NewFileOrchestrator(),
	}
}

func (factory *Factory) NewFeedOrchestrator() FeedOrchestrator {
	return FeedOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewFileOrchestrator() FileOrchestrator {
	return FileOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewNavigationOrchestrator() NavigationOrchestrator {
	return NavigationOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewOpenSearchDescriptionOrchestrator() OpenSearchDescriptionOrchestrator {
	return OpenSearchDescriptionOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewSearchOrchestrator() SearchOrchestrator {
	return SearchOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewSitemapOrchestrator() SitemapOrchestrator {
	return SitemapOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTagsOrchestrator() TagsOrchestrator {

	return TagsOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewViewModelOrchestrator() ViewModelOrchestrator {

	orchestrator := ViewModelOrchestrator{
		Orchestrator: factory.baseOrchestrator,

		navigationOrchestrator: factory.NewNavigationOrchestrator(),
		tagOrchestrator:        factory.NewTagsOrchestrator(),
		fileOrchestrator:       factory.NewFileOrchestrator(),
		locationOrchestrator:   factory.NewLocationOrchestrator(),
	}

	// refresh control
	go func() {
		for {
			select {
			case <-orchestrator.repository.AfterReindex():
				// reset the list
				factory.logger.Info("Resetting the leafes list")
				factory.baseOrchestrator.ResetCache()
			}
		}
	}()

	return orchestrator
}

func (factory *Factory) NewXmlSitemapOrchestrator() XmlSitemapOrchestrator {
	return XmlSitemapOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTypeAheadOrchestrator() TypeAheadOrchestrator {
	return TypeAheadOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTitlesOrchestrator() TitlesOrchestrator {
	return TitlesOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}

func (factory *Factory) NewUpdateOrchestrator() UpdateOrchestrator {
	return UpdateOrchestrator{
		Orchestrator:          factory.baseOrchestrator,
		viewModelOrchestrator: factory.NewViewModelOrchestrator(),
	}
}

func (factory *Factory) NewLocationOrchestrator() LocationOrchestrator {
	return LocationOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}
