// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeaheadhandler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator/typeaheadorchestrator"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel/typeaheadviewmodel"
	"io"
	"net/http"
)

func NewTitlesHandler(logger logger.Logger, config *config.Config, patherFactory paths.PatherFactory, itemIndex *index.Index) *TitleHandler {

	titlePathProvider := patherFactory.Absolute("/")
	titlesOrchestrator := typeaheadorchestrator.NewTitlesOrchestrator(itemIndex, titlePathProvider)

	return &TitleHandler{
		logger:             logger,
		itemIndex:          itemIndex,
		config:             config,
		patherFactory:      patherFactory,
		titlesOrchestrator: &titlesOrchestrator,
	}
}

type TitleHandler struct {
	logger             logger.Logger
	itemIndex          *index.Index
	config             *config.Config
	patherFactory      paths.PatherFactory
	titlesOrchestrator *typeaheadorchestrator.TitlesOrchestrator
}

func (handler *TitleHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the suggestions
		titles := handler.titlesOrchestrator.GetTitles()

		// set content type to json
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/json")

		// convert to json
		writeTitles(w, titles)
	}
}

func writeTitles(writer io.Writer, titles []typeaheadviewmodel.Title) error {
	bytes, err := json.MarshalIndent(titles, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
