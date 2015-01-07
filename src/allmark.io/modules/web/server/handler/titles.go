// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
)

type Titles struct {
	logger logger.Logger

	titlesOrchestrator *orchestrator.TitlesOrchestrator
}

func (handler *Titles) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "application/json; charset=utf-8")
		header.NoCache(w, r)
		header.VaryAcceptEncoding(w, r)

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
