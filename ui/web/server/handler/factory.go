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
	"github.com/andreaskoch/allmark2/ui/web/server/handler/debughandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/itemhandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/robotstxthandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/rsshandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/sitemaphandler"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/xmlsitemaphandler"
)

func NewRobotsTxtHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return robotstxthandler.New(logger, config, itemIndex, patherFactory)
}

func NewXmlSitemapHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return xmlsitemaphandler.New(logger, config, itemIndex, patherFactory)
}

func NewSitemapHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, patherFactory paths.PatherFactory) Handler {
	return sitemaphandler.New(logger, config, itemIndex, patherFactory)
}

func NewRssHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) Handler {
	return rsshandler.New(logger, config, itemIndex, fileIndex, patherFactory, converter)
}

func NewItemHandler(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) Handler {
	return itemhandler.New(logger, config, itemIndex, fileIndex, patherFactory, converter)
}

func NewDebugHandler(logger logger.Logger, itemIndex *index.ItemIndex, fileIndex *index.FileIndex) Handler {
	return debughandler.New(logger, itemIndex, fileIndex)
}
