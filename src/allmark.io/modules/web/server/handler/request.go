// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/route"
	"net/http"
)

func getRouteFromRequest(r *http.Request) route.Route {
	return route.NewFromRequest(r.URL.Path)
}

func getHostnameFromRequest(r *http.Request) string {
	return r.Host
}
