// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/debughandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/itemhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/opensearchdescriptionhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/robotstxthandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/rsshandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/searchhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/sitemaphandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/tagshandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/typeaheadhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/xmlsitemaphandler"
)

func NewErrorHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return errorhandler.New(logger, config, itemIndex, patherFactory)
}

func NewRobotsTxtHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return robotstxthandler.New(logger, config, itemIndex, patherFactory)
}

func NewXmlSitemapHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return xmlsitemaphandler.New(logger, config, itemIndex, patherFactory)
}

func NewTagsHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return tagshandler.New(logger, config, itemIndex, patherFactory)
}

func NewSitemapHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return sitemaphandler.New(logger, config, itemIndex, patherFactory)
}

func NewSearchHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex, searcher *search.ItemSearch) Handler {
	return searchhandler.New(logger, config, patherFactory, itemIndex, searcher)
}

func NewOpenSearchDescriptionHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex) Handler {
	return opensearchdescriptionhandler.New(logger, config, patherFactory, itemIndex)
}

func NewTypeAheadSearchHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex, searcher *search.ItemSearch) Handler {
	return typeaheadhandler.NewSearchHandler(logger, config, patherFactory, itemIndex, searcher)
}

func NewTypeAheadTitlesHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.ItemIndex) Handler {
	return typeaheadhandler.NewTitlesHandler(logger, config, patherFactory, itemIndex)
}

func NewRssHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory, converter conversion.Converter) Handler {
	return rsshandler.New(logger, config, itemIndex, patherFactory, converter)
}

func NewItemHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory, converter conversion.Converter) Handler {
	return itemhandler.New(logger, config, itemIndex, patherFactory, converter)
}

func NewDebugHandler(logger logger.Logger, itemIndex *index.ItemIndex) Handler {
	return debughandler.New(logger, itemIndex)
}
