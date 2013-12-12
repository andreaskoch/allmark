// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"github.com/andreaskoch/allmark2/ui/web/server/index"
	"net/http"
)

func New(logger logger.Logger, index *index.Index) *ItemHandler {
	return &ItemHandler{
		logger: logger,
		index:  index,
	}
}

type ItemHandler struct {
	logger logger.Logger
	index  *index.Index
}

func (handler *ItemHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the request route
		requestPath := handlerutil.GetRequestedPathFromRequest(r)
		requestRoute, err := route.NewFromRequest(requestPath)
		if err != nil {
			fmt.Fprintln(w, "%s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		item, found := handler.index.IsMatch(*requestRoute)
		if !found {
			fmt.Fprintln(w, "%s", "item not found")
			return
		}

		fmt.Fprintf(w, "%s", item.Content)
	}
}
