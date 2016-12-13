// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/andreaskoch/allmark/web/header"
	"github.com/andreaskoch/allmark/web/orchestrator"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"encoding/json"
	"io"
	"net/http"
)

func Titles(headerWriter header.HeaderWriter, titlesOrchestrator *orchestrator.TitlesOrchestrator) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// get the suggestions
		titles := titlesOrchestrator.GetTitles()
		writeTitles(w, titles)
	})

}

func writeTitles(writer io.Writer, titles []viewmodel.Title) error {
	bytes, err := json.MarshalIndent(titles, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
