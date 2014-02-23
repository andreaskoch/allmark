// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlerutil

import (
	"github.com/andreaskoch/allmark2/common/route"
	"net/http"
)

func GetRequestedPathFromRequest(r *http.Request) string {
	requestedRoute, err := route.NewFromRequest(r.URL.Path)
	if err != nil {
		panic(err)
	}

	requestedPath := requestedRoute.Value()
	return requestedPath
}

func GetHostnameFromRequest(r *http.Request) string {
	return r.Host
}
