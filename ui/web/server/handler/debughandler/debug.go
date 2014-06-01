// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debughandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"net/http"
)

func New(logger logger.Logger, itemIndex *index.Index) *DebugHandler {
	return &DebugHandler{
		logger:    logger,
		itemIndex: itemIndex,
	}
}

type DebugHandler struct {
	logger    logger.Logger
	itemIndex *index.Index
}

func (handler *DebugHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, handler.itemIndex.String())
	}
}
