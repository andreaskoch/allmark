// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"
)

func Redirect(baseUriTarget string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		redirectUrl := baseUriTarget + "/" + requestPath

		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
	})

}
