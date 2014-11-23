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

	updateChannel := make(chan bool, 1)
	repository.AfterReindex(updateChannel)

	// refresh control
	go func() {
		for {
			select {
			case <-updateChannel:
				// reset the list
				logger.Debug("Resetting the the cache")
				baseOrchestrator.ResetCache()
			}
		}
	}()

	return &Factory{
		logger: logger,

		baseOrchestrator: baseOrchestrator,
	}
}

type Factory struct {
	logger logger.Logger

	baseOrchestrator *Orchestrator

	viewModelOrchestrator             *ViewModelOrchestrator
	conversionModelOrchestrator       *ConversionModelOrchestrator
	feedOrchestrator                  *FeedOrchestrator
	fileOrchestrator                  *FileOrchestrator
	navigationOrchestrator            *NavigationOrchestrator
	openSearchDescriptionOrchestrator *OpenSearchDescriptionOrchestrator
	searchOrchestrator                *SearchOrchestrator
}

func (factory *Factory) NewConversionModelOrchestrator() *ConversionModelOrchestrator {

	if factory.conversionModelOrchestrator != nil {
		return factory.conversionModelOrchestrator
	}

	factory.conversionModelOrchestrator = &ConversionModelOrchestrator{
		Orchestrator:     factory.baseOrchestrator,
		fileOrchestrator: factory.NewFileOrchestrator(),
	}

	return factory.conversionModelOrchestrator
}

func (factory *Factory) NewFeedOrchestrator() *FeedOrchestrator {
	if factory.feedOrchestrator != nil {
		return factory.feedOrchestrator
	}

	factory.feedOrchestrator = &FeedOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.feedOrchestrator
}

func (factory *Factory) NewFileOrchestrator() *FileOrchestrator {

	if factory.fileOrchestrator != nil {
		return factory.fileOrchestrator
	}

	factory.fileOrchestrator = &FileOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.fileOrchestrator
}

func (factory *Factory) NewNavigationOrchestrator() *NavigationOrchestrator {

	if factory.navigationOrchestrator != nil {
		return factory.navigationOrchestrator
	}

	factory.navigationOrchestrator = &NavigationOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.navigationOrchestrator
}

func (factory *Factory) NewOpenSearchDescriptionOrchestrator() *OpenSearchDescriptionOrchestrator {

	if factory.openSearchDescriptionOrchestrator != nil {
		return factory.openSearchDescriptionOrchestrator
	}

	factory.openSearchDescriptionOrchestrator = &OpenSearchDescriptionOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.openSearchDescriptionOrchestrator
}

func (factory *Factory) NewSearchOrchestrator() *SearchOrchestrator {
	if factory.searchOrchestrator != nil {
		return factory.searchOrchestrator
	}

	factory.searchOrchestrator = &SearchOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.searchOrchestrator
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

func (factory *Factory) NewViewModelOrchestrator() *ViewModelOrchestrator {

	// cache lookup
	if factory.viewModelOrchestrator != nil {
		return factory.viewModelOrchestrator
	}

	orchestrator := &ViewModelOrchestrator{
		Orchestrator: factory.baseOrchestrator,

		navigationOrchestrator: factory.NewNavigationOrchestrator(),
		tagOrchestrator:        factory.NewTagsOrchestrator(),
		fileOrchestrator:       factory.NewFileOrchestrator(),
		locationOrchestrator:   factory.NewLocationOrchestrator(),
	}

	// warm up the caches
	orchestrator.blockingCacheWarmup()

	// store
	factory.viewModelOrchestrator = orchestrator

	return factory.viewModelOrchestrator
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
