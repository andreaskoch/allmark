// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"io"
	"net/http"
)

type Titles struct {
	logger logger.Logger

	titlesOrchestrator orchestrator.TitlesOrchestrator
}

func (handler *Titles) Func() func(w http.ResponseWriter, r *http.Request) {

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

func writeTitles(writer io.Writer, titles []viewmodel.Title) error {
	bytes, err := json.MarshalIndent(titles, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
