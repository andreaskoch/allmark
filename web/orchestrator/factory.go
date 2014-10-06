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
		factory.baseOrchestrator,
		factory.NewFileOrchestrator(),
	}
}

func (factory *Factory) NewFeedOrchestrator() FeedOrchestrator {
	return FeedOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewFileOrchestrator() FileOrchestrator {
	return FileOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewNavigationOrchestrator() NavigationOrchestrator {
	return NavigationOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewOpenSearchDescriptionOrchestrator() OpenSearchDescriptionOrchestrator {
	return OpenSearchDescriptionOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewSearchOrchestrator() SearchOrchestrator {
	return SearchOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewSitemapOrchestrator() SitemapOrchestrator {
	return SitemapOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTagsOrchestrator() TagsOrchestrator {

	return TagsOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewViewModelOrchestrator() ViewModelOrchestrator {

	orchestrator := ViewModelOrchestrator{
		factory.baseOrchestrator,

		factory.NewNavigationOrchestrator(),
		factory.NewTagsOrchestrator(),
		factory.NewFileOrchestrator(),
		factory.NewLocationOrchestrator(),
	}

	// refresh control
	go func() {
		for {
			select {
			case <-orchestrator.repository.AfterReindex():
				// reset the list
				factory.logger.Info("Resetting the leafes list")
				orchestrator.ResetCache()
			}
		}
	}()

	return orchestrator
}

func (factory *Factory) NewXmlSitemapOrchestrator() XmlSitemapOrchestrator {
	return XmlSitemapOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTypeAheadOrchestrator() TypeAheadOrchestrator {
	return TypeAheadOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewTitlesOrchestrator() TitlesOrchestrator {
	return TitlesOrchestrator{
		factory.baseOrchestrator,
	}
}

func (factory *Factory) NewUpdateOrchestrator() UpdateOrchestrator {
	return UpdateOrchestrator{
		factory.baseOrchestrator,
		factory.NewViewModelOrchestrator(),
	}
}

func (factory *Factory) NewLocationOrchestrator() LocationOrchestrator {
	return LocationOrchestrator{
		factory.baseOrchestrator,
		nil,
	}
}
