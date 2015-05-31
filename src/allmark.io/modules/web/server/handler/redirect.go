// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"net/http"
)

type Redirect struct {
	logger        logger.Logger
	baseUriTarget string
}

func (handler *Redirect) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		redirectUrl := handler.baseUriTarget + "/" + requestPath

		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
	}

}
