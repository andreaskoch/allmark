// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"net/http"
)

func Redirect(logger logger.Logger, baseURITarget string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestPath := r.URL.Path
		redirectURL := baseURITarget + "/" + requestPath

		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	})

}
