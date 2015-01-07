// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/web/server/header"
	"fmt"
	"net/http"
)

type RobotsTxt struct {
	logger logger.Logger
}

func (handler *RobotsTxt) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// set headers
		header.ContentType(w, r, "text/plain; charset=utf-8")
		header.Cache(w, r, header.STATICCONTENT_CACHEDURATION_SECONDS)
		header.VaryAcceptEncoding(w, r)

		fmt.Fprintf(w, `User-agent: *
Disallow: /thumbnails
Disallow: /*.rtf$
Disallow: /*.json$
Disallow: /*.print$
Disallow: /*.ws$

Sitemap: http://%s/sitemap.xml`, getHostnameFromRequest(r))
	}
}
