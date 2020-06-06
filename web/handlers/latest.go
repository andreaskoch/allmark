// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/util/hashutil"
	"github.com/elWyatt/allmark/web/header"
	"github.com/elWyatt/allmark/web/orchestrator"
	"encoding/json"
	"net/http"
	"strings"
)

func Latest(logger logger.Logger, headerWriter header.HeaderWriter, viewModelOrchestrator *orchestrator.ViewModelOrchestrator, fallbackHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set headers
		headerWriter.Write(w, header.CONTENTTYPE_JSON)

		// strip the "latest" or ".latest" suffix from the path
		path := r.URL.Path
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
