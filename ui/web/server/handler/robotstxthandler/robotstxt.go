// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package robotstxthandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/handlerutil"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory) *RobotsTxtHandler {

	return &RobotsTxtHandler{
		logger:        logger,
		itemIndex:     itemIndex,
		config:        config,
		patherFactory: patherFactory,
	}
}

type RobotsTxtHandler struct {
	logger        logger.Logger
	itemIndex     *index.Index
	config        *config.Config
	patherFactory paths.PatherFactory
}

func (handler *RobotsTxtHandler) Func() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `User-agent: *
allow: /

Sitemap: http://%s/sitemap.xml`, handlerutil.GetHostnameFromRequest(r))
	}
}
