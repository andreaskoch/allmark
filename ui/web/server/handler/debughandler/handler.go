// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debughandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/ui/web/server/index"
	"net/http"
)

func New(logger logger.Logger, index *index.Index) *DebugHandler {
	return &DebugHandler{
		logger: logger,
		index:  index,
	}
}

type DebugHandler struct {
	logger logger.Logger
	index  *index.Index
}

func (handler *DebugHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, route := range handler.index.Routes() {
			fmt.Fprintf(w, "%q\n", route.Value())
		}
	}
}
