// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"net/http"
)

type RobotsTxt struct {
	logger logger.Logger
}

func (handler *RobotsTxt) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `User-agent: *
allow: /

Sitemap: http://%s/sitemap.xml`, getHostnameFromRequest(r))
	}
}
