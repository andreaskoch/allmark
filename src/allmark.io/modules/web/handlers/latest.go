// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/route"
	"allmark.io/modules/common/util/hashutil"
	"allmark.io/modules/web/orchestrator"
	"allmark.io/modules/web/header"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func Latest(logger logger.Logger, headerWriter header.HeaderWriter, viewModelOrchestrator *orchestrator.ViewModelOrchestrator, fallbackHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// strip the "latest" or ".latest" suffix from the path
		path = strings.TrimSuffix(path, "latest")
		path = strings.TrimSuffix(path, ".")

		// get the request route
		requestRoute := route.NewFromRequest(path)

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if latestModels, found := viewModelOrchestrator.GetLatest(requestRoute, 3, 1); found {

			// convert the viewmodel to json
			jsonBytes, err := json.MarshalIndent(latestModels, "", "\t")
			if err != nil {
				logger.Error("Unable to convert viewmodel to json. Error: %s", err)
				return
			}

			// etag cache validator
			etag := hashutil.FromBytes(jsonBytes)
			if etag != "" {
				header.ETag(w, etag)
			}

			w.Write(jsonBytes)

			return
		}

		// fallback
		fallbackHandler.ServeHTTP(w, r)
	})
}
