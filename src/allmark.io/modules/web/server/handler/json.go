// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/server/header"
	"allmark.io/modules/web/view/viewmodel"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Json struct {
	logger logger.Logger

	viewModelOrchestrator *orchestrator.ViewModelOrchestrator

	fallbackHandler Handler
}

func (handler *Json) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "application/json; charset=utf-8")
		header.Cache(w, r, header.DYNAMICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "json" or ".json" suffix from the path
		path = strings.TrimSuffix(path, "json")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if viewModel, found := handler.viewModelOrchestrator.GetFullViewModel(requestRoute); found {
			writeViewModelAsJson(w, viewModel)
			return
		}

		// fallback to the item handler
		fallbackHandler := handler.fallbackHandler.Func()
		fallbackHandler(w, r)
	}
}

func writeViewModelAsJson(writer io.Writer, viewModel viewmodel.Model) error {
	bytes, err := json.MarshalIndent(viewModel, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
