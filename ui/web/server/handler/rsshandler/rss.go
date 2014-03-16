// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsshandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	// "github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	// "github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/view/templates"
	// "github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	// "io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.ItemIndex, fileIndex *index.FileIndex, patherFactory paths.PatherFactory, converter conversion.Converter) *RssHandler {

	templateProvider := templates.NewProvider(".")
	viewModelOrchestrator := orchestrator.NewViewModelOrchestrator(itemIndex, converter)

	return &RssHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		fileIndex:             fileIndex,
		config:                config,
		patherFactory:         patherFactory,
		templateProvider:      templateProvider,
		viewModelOrchestrator: viewModelOrchestrator,
	}
}

type RssHandler struct {
	logger                logger.Logger
	itemIndex             *index.ItemIndex
	fileIndex             *index.FileIndex
	config                *config.Config
	patherFactory         paths.PatherFactory
	templateProvider      *templates.Provider
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
}

func (handler *RssHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// read the page url-parameter
		page, pageParameterIsAvailable := handlerutil.GetPageParameterFromUrl(*r.URL)
		if !pageParameterIsAvailable || page == 0 {
			page = 1
		}

		fmt.Fprintf(w, "%s", page)
	}
}
