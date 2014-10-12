// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/server/header"
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

		// set headers
		header.ContentType(w, r, "application/json")
		header.NoCache(w, r)

		// get the suggestions
		titles := handler.titlesOrchestrator.GetTitles()
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
