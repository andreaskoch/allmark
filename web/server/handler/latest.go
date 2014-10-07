// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/orchestrator"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Latest struct {
	logger logger.Logger

	viewModelOrchestrator orchestrator.ViewModelOrchestrator
	fallbackHandler       Handler
}

func (handler *Latest) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "latest" or ".latest" suffix from the path
		path = strings.TrimSuffix(path, "latest")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if latestModels, found := handler.viewModelOrchestrator.GetLatest(requestRoute, 3, 1); found {
			writeViewModelsAsJson(w, latestModels)
			return
		}

		// fallback to the item handler
		fallbackHandler := handler.fallbackHandler.Func()
		fallbackHandler(w, r)
	}
}

func writeViewModelsAsJson(writer io.Writer, viewModels []*viewmodel.Model) error {
	bytes, err := json.MarshalIndent(viewModels, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
