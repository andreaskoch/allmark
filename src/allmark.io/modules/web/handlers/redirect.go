// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"allmark.io/modules/common/logger"
	"net/http"
)

func Redirect(logger logger.Logger, baseURITarget string, baseHandler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestPath := r.URL.Path
		redirectURL := baseURITarget + "/" + requestPath

		// don't redirect if it's a local request
		if isLocalRequest(r) {
			logger.Debug("Skipping the redirect to %q for request %q because it's a local request", redirectURL, r.URL.String())
			baseHandler.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	})

}
