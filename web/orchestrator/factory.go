// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/elWyatt/allmark/common/config"
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/dataaccess"
	"github.com/elWyatt/allmark/services/converter"
	"github.com/elWyatt/allmark/services/parser"
	"github.com/elWyatt/allmark/web/webpaths"
)

func NewFactory(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Factory {

	baseOrchestrator := newBaseOrchestrator(logger, config, repository, parser, converter, webPathProvider)

	// listen for updates
	repositoryUpdates := make(chan dataaccess.Update, 1)
	repository.Subscribe(repositoryUpdates)

	go func() {
		for update := range repositoryUpdates {
			logger.Info("Received and update (%s). Resetting the the cache.", update.String())
			baseOrchestrator.UpdateCache(update)
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
	sitemapOrchestrator               *SitemapOrchestrator
	tagsOrchestrator                  *TagsOrchestrator
	xmlSitemapOrchestrator            *XmlSitemapOrchestrator
	typeAheadOrchestrator             *TypeAheadOrchestrator
	titlesOrchestrator                *TitlesOrchestrator
	updateOrchestrator                *UpdateOrchestrator
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

func (factory *Factory) NewSitemapOrchestrator() *SitemapOrchestrator {

	if factory.sitemapOrchestrator != nil {
		return factory.sitemapOrchestrator
	}

	factory.sitemapOrchestrator = &SitemapOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.sitemapOrchestrator
}

func (factory *Factory) NewTagsOrchestrator() *TagsOrchestrator {

	if factory.tagsOrchestrator != nil {
		return factory.tagsOrchestrator
	}

	factory.tagsOrchestrator = &TagsOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.tagsOrchestrator
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
	}

	// store
	factory.viewModelOrchestrator = orchestrator

	return factory.viewModelOrchestrator
}

func (factory *Factory) NewXMLSitemapOrchestrator() *XmlSitemapOrchestrator {

	if factory.xmlSitemapOrchestrator != nil {
		return factory.xmlSitemapOrchestrator
	}

	factory.xmlSitemapOrchestrator = &XmlSitemapOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.xmlSitemapOrchestrator
}

func (factory *Factory) NewTypeAheadOrchestrator() *TypeAheadOrchestrator {

	if factory.typeAheadOrchestrator != nil {
		return factory.typeAheadOrchestrator
	}

	factory.typeAheadOrchestrator = &TypeAheadOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.typeAheadOrchestrator
}

func (factory *Factory) NewTitlesOrchestrator() *TitlesOrchestrator {

	if factory.titlesOrchestrator != nil {
		return factory.titlesOrchestrator
	}

	factory.titlesOrchestrator = &TitlesOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}

	return factory.titlesOrchestrator
}

func (factory *Factory) NewUpdateOrchestrator() *UpdateOrchestrator {
	if factory.updateOrchestrator != nil {
		return factory.updateOrchestrator
	}

	factory.updateOrchestrator = &UpdateOrchestrator{
		Orchestrator:          factory.baseOrchestrator,
		viewModelOrchestrator: factory.NewViewModelOrchestrator(),
	}

	return factory.updateOrchestrator
}

// NewAliasIndexOrchestrator creates a new alias-index orchestrator.
func (factory *Factory) NewAliasIndexOrchestrator() *AliasIndexOrchestrator {
	return &AliasIndexOrchestrator{
		Orchestrator: factory.baseOrchestrator,
	}
}
